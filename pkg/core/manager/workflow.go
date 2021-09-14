package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/fuseml/fuseml-core/pkg/domain"
)

// createWorkflowListenerTimeout is the time (in minutes) that FuseML waits for the workflow listener
// to be available
const createWorkflowListenerTimeout = 1

// WorkflowManager implements the domain.WorkflowManager interface
type WorkflowManager struct {
	workflowBackend   domain.WorkflowBackend
	workflowStore     domain.WorkflowStore
	codesetStore      domain.CodesetStore
	extensionRegistry domain.ExtensionRegistry
}

// NewWorkflowManager initializes a Workflow Manager
// FIXME: instead of CodesetStore, receive a CodesetManager
func NewWorkflowManager(
	workflowBackend domain.WorkflowBackend,
	workflowStore domain.WorkflowStore,
	codesetStore domain.CodesetStore,
	extensionRegistry domain.ExtensionRegistry) *WorkflowManager {
	return &WorkflowManager{workflowBackend, workflowStore, codesetStore, extensionRegistry}
}

// GetWorkflows returns a list of Workflows.
func (mgr *WorkflowManager) GetWorkflows(ctx context.Context, name *string) []*domain.Workflow {
	return mgr.workflowStore.GetWorkflows(ctx, name)
}

// CreateWorkflow creates a new Workflow.
func (mgr *WorkflowManager) CreateWorkflow(ctx context.Context, wf *domain.Workflow) (*domain.Workflow, error) {
	wf.Created = time.Now()
	err := mgr.resolveExtensionReferences(ctx, wf)
	if err != nil {
		return nil, err
	}
	err = mgr.workflowBackend.CreateWorkflow(ctx, wf)
	if err != nil {
		return nil, err
	}
	return mgr.workflowStore.AddWorkflow(ctx, wf)
}

// GetWorkflow retrieves a Workflow.
func (mgr *WorkflowManager) GetWorkflow(ctx context.Context, name string) (*domain.Workflow, error) {
	return mgr.workflowStore.GetWorkflow(ctx, name)
}

// DeleteWorkflow deletes a Workflow and its assignments.
func (mgr *WorkflowManager) DeleteWorkflow(ctx context.Context, name string) error {
	// unassign all assigned codesets, if there's any
	codesetAssignments := mgr.workflowStore.GetCodesetAssignments(ctx, name)
	for _, ca := range codesetAssignments {
		err := mgr.UnassignFromCodeset(ctx, name, ca.Codeset.Project, ca.Codeset.Name)
		if err != nil {
			return err
		}
	}

	// delete tekton pipeline
	err := mgr.workflowBackend.DeleteWorkflow(ctx, name)
	if err != nil {
		return err
	}

	// delete workflow
	err = mgr.workflowStore.DeleteWorkflow(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// AssignToCodeset assigns a Workflow to a Codeset.
func (mgr *WorkflowManager) AssignToCodeset(ctx context.Context, name, codesetProject, codesetName string) (wfListener *domain.WorkflowListener, webhookID *int64, err error) {
	_, err = mgr.workflowStore.GetWorkflow(ctx, name)
	if err != nil {
		return nil, nil, err
	}

	codeset, err := mgr.codesetStore.Find(ctx, codesetProject, codesetName)
	if err != nil {
		return nil, nil, err
	}

	wfListener, err = mgr.workflowBackend.CreateWorkflowListener(ctx, name, createWorkflowListenerTimeout*time.Minute)
	if err != nil {
		return nil, nil, err
	}

	assignment, err := mgr.workflowStore.GetCodesetAssignment(ctx, name, codeset)
	if err == nil {
		return wfListener, assignment.WebhookID, nil
	}

	webhookID, err = mgr.codesetStore.CreateWebhook(ctx, codeset, wfListener.URL)
	if err != nil {
		return nil, nil, err
	}

	mgr.workflowStore.AddCodesetAssignment(ctx, name, codeset, webhookID)
	mgr.codesetStore.Subscribe(ctx, mgr, codeset)
	mgr.workflowBackend.CreateWorkflowRun(ctx, name, codeset)
	return
}

// UnassignFromCodeset unassign a Workflow from a Codeset
func (mgr *WorkflowManager) UnassignFromCodeset(ctx context.Context, name, codesetProject, codesetName string) (err error) {
	codeset, err := mgr.codesetStore.Find(ctx, codesetProject, codesetName)
	if err != nil {
		return err
	}

	assignment, err := mgr.workflowStore.GetCodesetAssignment(ctx, name, codeset)
	if err != nil {
		return err
	}

	if assignment.WebhookID != nil {
		err = mgr.codesetStore.DeleteWebhook(ctx, codeset, assignment.WebhookID)
		if err != nil {
			return err
		}
	}

	if len(mgr.workflowStore.GetCodesetAssignments(ctx, name)) == 1 {
		err = mgr.workflowBackend.DeleteWorkflowListener(ctx, name)
		if err != nil {
			return err
		}
	}

	mgr.workflowStore.DeleteCodesetAssignment(ctx, name, codeset)
	mgr.codesetStore.Unsubscribe(ctx, mgr, codeset)
	return
}

// GetAllCodesetAssignments lists Workflow assignments.
func (mgr *WorkflowManager) GetAllCodesetAssignments(ctx context.Context, name *string) map[string][]*domain.CodesetAssignment {
	return mgr.workflowStore.GetAllCodesetAssignments(ctx, name)
}

// GetAssignmentStatus returns the status of a Workflow assignment.
func (mgr *WorkflowManager) GetAssignmentStatus(ctx context.Context, name string) *domain.WorkflowAssignmentStatus {
	status := domain.WorkflowAssignmentStatus{}
	listener, err := mgr.workflowBackend.GetWorkflowListener(ctx, name)
	if err != nil {
		return &status
	}

	status.Available = listener.Available
	status.URL = listener.DashboardURL
	return &status
}

// GetWorkflowRuns returns a lists Workflow runs.
func (mgr *WorkflowManager) GetWorkflowRuns(ctx context.Context, filter *domain.WorkflowRunFilter) ([]*domain.WorkflowRun, error) {
	workflowRuns := []*domain.WorkflowRun{}
	var wfName *string
	if filter != nil {
		wfName = filter.WorkflowName
	}
	workflows := mgr.workflowStore.GetWorkflows(ctx, wfName)

	for _, workflow := range workflows {
		runs, err := mgr.workflowBackend.GetWorkflowRuns(ctx, workflow, filter)
		if err != nil {
			return nil, err
		}
		workflowRuns = append(workflowRuns, runs...)
	}

	return workflowRuns, nil
}

// OnDeletingCodeset perform operations on workflows when a codeset is deleted
func (mgr *WorkflowManager) OnDeletingCodeset(ctx context.Context, codeset *domain.Codeset) {
	for _, wf := range mgr.GetWorkflows(ctx, nil) {
		mgr.UnassignFromCodeset(ctx, wf.Name, codeset.Project, codeset.Name)
	}
}

// Resolve all the extension references in the workflow steps and update them with actual
// extension endpoints and credentials
func (mgr *WorkflowManager) resolveExtensionReferences(ctx context.Context, wf *domain.Workflow) error {
	for _, step := range wf.Steps {
		for _, extReq := range step.Extensions {
			accessDescList, err := mgr.extensionRegistry.GetExtensionAccessDescriptors(ctx, &domain.ExtensionQuery{
				ExtensionID:        extReq.ExtensionID,
				Product:            extReq.Product,
				VersionConstraints: extReq.VersionConstraints,
				Zone:               extReq.Zone,
				// allow extensions outside of the zone for now
				StrictZoneMatch: false,
				ServiceID:       extReq.ServiceID,
				ServiceResource: extReq.ServiceResource,
				ServiceCategory: extReq.ServiceCategory,
				// determine endpoint type automatically based on zone
				Type: nil,
				// only global credentials supported for now
				CredentialsScope: domain.ECSGlobal,
			})
			if err != nil {
				return fmt.Errorf("error resolving extension requirements for step %q extension %q: %w", step.Name, extReq.Name, err)
			}
			if len(accessDescList) == 0 {
				return fmt.Errorf("could not resolve extension requirements for step %q extension %q", step.Name, extReq.Name)
			}
			// for now, assume that all internal endpoints are accessible from workflow steps and
			// prefer internal endpoints if more results are returned
			extReq.ExtensionAccess = accessDescList[0]
			for _, accessDesc := range accessDescList {
				if accessDesc.Endpoint.Type == domain.EETInternal {
					extReq.ExtensionAccess = accessDesc
					break
				}
			}
		}
	}

	return nil
}
