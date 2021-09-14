package manager

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tektoncd/pipeline/test/diff"

	"github.com/fuseml/fuseml-core/pkg/core"
	"github.com/fuseml/fuseml-core/pkg/domain"
	"github.com/fuseml/fuseml-core/pkg/util"
)

const errCodesetNotFound = codesetErr("codeset not found")

var (
	// workflowBackend stores WorkflowListener and WorkflowRuns for a workflow created by fakeWorkflowBackend
	workflowBackend domain.WorkflowBackend

	// workflowStore stores Workflow and Assignments
	workflowStore domain.WorkflowStore

	// codesetStore stores codesets that are created when initializing fakeWorkflowManager
	// The following codesets are created when calling newFakeWorkflowManager:
	// 1. name: cs0, project: csproject0
	// 2. name: cs1, project: csproject1
	// 3. name: cs2, project: csproject1
	codesetStore *fakeCodesetStore

	// extensionRegistry stores extensions
	extensionRegistry *ExtensionRegistry

	// workflowRunStatuses are the possible Status for a WorkflowRun. The status of a WorkflowRun is set
	// accordingly to its order, cycling between the workflowRunStatuses. E.g. run0: Succeeded, run1: Failed,
	// run2: Succeeded, ...
	workflowRunStatuses = []string{"Succeeded", "Failed"}
)

type codesetErr string

func (e codesetErr) Error() string {
	return string(e)
}

func TestCreateWorkflow(t *testing.T) {
	t.Run("new workflow", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		wf := domain.Workflow{Name: "test"}
		got, err := mgr.CreateWorkflow(context.Background(), &wf)
		assertError(t, err, nil)

		want, _ := workflowStore.GetWorkflow(context.TODO(), wf.Name)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow: %s", diff.PrintWantGot(d))
		}

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		err = workflowBackend.CreateWorkflowRun(context.TODO(), wf.Name, codesets[0])
		assertError(t, err, nil)
	})

	t.Run("existing workflow", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		wf := domain.Workflow{Name: "test"}
		_, err := mgr.CreateWorkflow(context.Background(), &wf)
		assertError(t, err, nil)

		_, err = mgr.CreateWorkflow(context.Background(), &wf)
		assertError(t, err, domain.ErrWorkflowExists)

		got := workflowStore.GetWorkflows(context.TODO(), nil)
		want := []*domain.Workflow{&wf}
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow: %s", diff.PrintWantGot(d))
		}
	})

	t.Run("new workflow with matching extensions", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		ext := createFakeExtension(t, mgr, "test-")
		ext, err := mgr.extensionRegistry.RegisterExtension(context.Background(), ext)
		assertError(t, err, nil)

		wf := domain.Workflow{
			Name: "test",
			Steps: []*domain.WorkflowStep{{
				Name: "test-step",
				Extensions: []*domain.WorkflowStepExtension{{
					Name:               "test-extension",
					Product:            ext.Product,
					VersionConstraints: ">=" + ext.Version,
					ServiceResource:    ext.ListServices()[0].Resource,
				}},
			}},
		}
		got, err := mgr.CreateWorkflow(context.Background(), &wf)
		assertError(t, err, nil)

		want, _ := workflowStore.GetWorkflow(context.TODO(), wf.Name)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow: %s", diff.PrintWantGot(d))
		}

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		err = workflowBackend.CreateWorkflowRun(context.TODO(), wf.Name, codesets[0])
		assertError(t, err, nil)
	})

	t.Run("new workflow with no matching extensions", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		ext := createFakeExtension(t, mgr, "test-")
		ext, err := mgr.extensionRegistry.RegisterExtension(context.Background(), ext)
		assertError(t, err, nil)

		wf := domain.Workflow{
			Name: "test",
			Steps: []*domain.WorkflowStep{{
				Name: "test-step",
				Extensions: []*domain.WorkflowStepExtension{{
					Name:               "test-extension",
					Product:            ext.Product,
					VersionConstraints: "<" + ext.Version,
					ServiceResource:    ext.ListServices()[0].Resource,
				}},
			}},
		}
		_, err = mgr.CreateWorkflow(context.Background(), &wf)
		if err.Error() != "could not resolve extension requirements for step \"test-step\" extension \"test-extension\"" {
			t.Errorf("Unexpected Workflow error: %q", err)
		}
	})
}

func TestGetWorkflows(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		want := []*domain.Workflow{}

		// no workflows
		got := mgr.GetWorkflows(context.TODO(), nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow list: %s", diff.PrintWantGot(d))
		}

		// create 3 workflows (wf0, wf1, wf2)
		for i := 0; i < 3; i++ {
			wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: fmt.Sprintf("wf%d", i)})
			assertError(t, err, nil)
			want = append(want, wf)
		}

		got = mgr.GetWorkflows(context.TODO(), nil)
		if d := cmp.Diff(want, got, cmpopts.SortSlices(func(x, y *domain.Workflow) bool { return x.Name < y.Name })); d != "" {
			t.Errorf("Unexpected Workflow list: %s", diff.PrintWantGot(d))
		}
	})

	t.Run("by workflow name", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		want := []*domain.Workflow{}

		// no workflows
		wfName := "does-not-exist"
		got := mgr.GetWorkflows(context.TODO(), &wfName)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow list: %s", diff.PrintWantGot(d))
		}

		// create 3 workflows (wf0, wf1, wf2)
		for i := 0; i < 3; i++ {
			wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: fmt.Sprintf("wf%d", i)})
			assertError(t, err, nil)
			want = append(want, wf)
		}

		for i := 0; i < len(want); i++ {
			name := fmt.Sprintf("wf%d", i)
			got := mgr.GetWorkflows(context.TODO(), &name)
			if d := cmp.Diff([]*domain.Workflow{want[i]}, got, cmpopts.SortSlices(func(x, y *domain.Workflow) bool { return x.Name < y.Name })); d != "" {
				t.Errorf("Unexpected Workflow list: %s", diff.PrintWantGot(d))
			}
		}
	})
}

func TestGetWorkflow(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		want, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		got, err := mgr.GetWorkflow(context.Background(), want.Name)
		assertError(t, err, nil)

		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow: %s", diff.PrintWantGot(d))
		}
	})

	t.Run("not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		_, err := mgr.GetWorkflow(context.Background(), "wf")
		assertError(t, err, domain.ErrWorkflowNotFound)
	})
}

func TestDeleteWorkflow(t *testing.T) {
	t.Run("not assigned", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		err = mgr.DeleteWorkflow(context.Background(), wf.Name)
		assertError(t, err, nil)

		got := workflowStore.GetWorkflows(context.TODO(), &wf.Name)
		if d := cmp.Diff([]*domain.Workflow{}, got); d != "" {
			t.Errorf("Unexpected Workflow: %s", diff.PrintWantGot(d))
		}

		err = workflowBackend.CreateWorkflowRun(context.TODO(), wf.Name, nil)
		assertStrings(t, err.Error(), "workflow not found")

	})

	t.Run("assigned", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		_, _, got := mgr.AssignToCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
		assertError(t, got, nil)

		err = mgr.DeleteWorkflow(context.Background(), wf.Name)
		assertError(t, err, nil)
	})

	t.Run("not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)
		err := mgr.DeleteWorkflow(context.Background(), "wf")
		assertError(t, err, nil)
	})
}

func TestAssignToCodeset(t *testing.T) {
	t.Run("assign", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		codeset := codesets[0]
		wantListener, webhookID, err := mgr.AssignToCodeset(context.Background(), wf.Name, codeset.Project, codeset.Name)
		assertError(t, err, nil)

		ignoreUnexported := cmpopts.IgnoreUnexported(WorkflowManager{})
		gotSubscribers := codesetStore.getSubscribers(context.TODO(), codeset)
		if d := cmp.Diff([]domain.CodesetSubscriber{mgr}, gotSubscribers, ignoreUnexported); d != "" {
			t.Errorf("Unexpected codeset subscriber: %s", diff.PrintWantGot(d))
		}

		got := workflowStore.GetAllCodesetAssignments(context.TODO(), &wf.Name)
		csAsg := domain.CodesetAssignment{Codeset: codeset, WebhookID: webhookID}
		want := map[string][]*domain.CodesetAssignment{wf.Name: {&csAsg}}
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		gotListener, err := workflowBackend.GetWorkflowListener(context.TODO(), wf.Name)
		assertError(t, err, nil)
		if d := cmp.Diff(wantListener, gotListener); d != "" {
			t.Errorf("Unexpected Listener: %s", diff.PrintWantGot(d))
		}

		workflowRuns, err := workflowBackend.GetWorkflowRuns(context.TODO(), wf, nil)
		assertError(t, err, nil)
		gotRuns := len(workflowRuns)
		wantRuns := 1
		if gotRuns != wantRuns {
			t.Errorf("Expected %d WorkflowRun got %d", wantRuns, gotRuns)
		}
	})

	t.Run("existing assignment", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)

		for i := 0; i < 2; i++ {
			_, _, err := mgr.AssignToCodeset(context.TODO(), wf.Name, codesets[0].Project, codesets[0].Name)
			assertError(t, err, nil)
		}

		_, err = workflowBackend.GetWorkflowListener(context.TODO(), wf.Name)
		assertError(t, err, nil)

		workflowRuns, err := workflowBackend.GetWorkflowRuns(context.TODO(), wf, nil)
		assertError(t, err, nil)

		gotRuns := len(workflowRuns)
		wantRuns := 1
		if gotRuns != wantRuns {
			t.Errorf("Expected %d WorkflowRun got %d", wantRuns, gotRuns)
		}
	})

	t.Run("workflow not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wfName := "unknownWf"
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		_, _, got := mgr.AssignToCodeset(context.Background(), wfName, codesets[0].Project, codesets[0].Name)
		assertError(t, got, domain.ErrWorkflowNotFound)

		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), nil)
		wantAss := map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		_, err := workflowBackend.GetWorkflowListener(context.TODO(), wfName)
		assertStrings(t, err.Error(), "listener not found")

	})

	t.Run("codeset not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		_, _, got := mgr.AssignToCodeset(context.Background(), wf.Name, "unknownProj", "unknownCs")
		assertError(t, got, errCodesetNotFound)

		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), nil)
		wantAss := map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		_, err = workflowBackend.GetWorkflowListener(context.TODO(), wf.Name)
		assertStrings(t, err.Error(), "listener not found")
	})
}

func TestUnassignFromCodeset(t *testing.T) {
	t.Run("unassign", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		var listener *domain.WorkflowListener
		var webhookID *int64
		webhooks := map[*domain.Codeset][]*int64{}
		for i := 0; i < 2; i++ {
			codeset := codesets[i]
			listener, webhookID, err = mgr.AssignToCodeset(context.Background(), wf.Name, codeset.Project, codeset.Name)
			assertError(t, err, nil)

			if webhook, exists := webhooks[codeset]; exists {
				webhooks[codeset] = append(webhook, webhookID)
			} else {
				webhooks[codeset] = []*int64{webhookID}
			}
		}

		// delete wf assignment to cs0
		err = mgr.UnassignFromCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
		assertError(t, err, nil)
		gotSubscribers := codesetStore.getSubscribers(context.TODO(), codesets[0])
		if d := cmp.Diff([]domain.CodesetSubscriber{}, gotSubscribers); d != "" {
			t.Errorf("Unexpected codeset subscriber: %s", diff.PrintWantGot(d))
		}

		// should have only one assignment to cs1
		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), &wf.Name)
		csAsg := domain.CodesetAssignment{Codeset: codesets[1], WebhookID: webhooks[codesets[1]][0]}
		wantAss := map[string][]*domain.CodesetAssignment{wf.Name: {&csAsg}}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		// listener should still exist
		gotListener, err := workflowBackend.GetWorkflowListener(context.TODO(), wf.Name)
		assertError(t, err, nil)
		if d := cmp.Diff(listener, gotListener); d != "" {
			t.Errorf("Unexpected Listener: %s", diff.PrintWantGot(d))
		}

		// delete wf assignment to cs1
		err = mgr.UnassignFromCodeset(context.Background(), wf.Name, codesets[1].Project, codesets[1].Name)
		assertError(t, err, nil)

		// should have no assignment
		gotAss = workflowStore.GetAllCodesetAssignments(context.TODO(), &wf.Name)
		wantAss = map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		// listener should be gone
		_, err = workflowBackend.GetWorkflowListener(context.TODO(), wf.Name)
		assertStrings(t, err.Error(), "listener not found")
	})

	t.Run("workflow not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wfName := "unknownWf"
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		got := mgr.UnassignFromCodeset(context.Background(), wfName, codesets[0].Project, codesets[0].Name)
		assertError(t, got, domain.ErrWorkflowNotFound)

		// should have no assignment
		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), nil)
		wantAss := map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		// should have no listener
		_, err := workflowBackend.GetWorkflowListener(context.TODO(), wfName)
		assertStrings(t, err.Error(), "listener not found")
	})

	t.Run("codeset not found", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wfName := "wf"
		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: wfName})
		assertError(t, err, nil)

		got := mgr.UnassignFromCodeset(context.Background(), wf.Name, "unknownProj", "unknownCs")
		assertError(t, got, errCodesetNotFound)

		// should have no assignment
		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), nil)
		wantAss := map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}

		// should have no listener
		_, err = workflowBackend.GetWorkflowListener(context.TODO(), wfName)
		assertStrings(t, err.Error(), "listener not found")
	})

	t.Run("workflow not assigned", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		got := mgr.UnassignFromCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
		assertError(t, got, domain.ErrWorkflowNotAssignedToCodeset)
	})

	t.Run("on codeset deleting", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		codeset := codesets[0]
		_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codeset.Project, codeset.Name)
		assertError(t, err, nil)

		codesetStore.Delete(context.TODO(), codeset.Project, codeset.Name)

		// should have no assignment
		gotAss := workflowStore.GetAllCodesetAssignments(context.TODO(), nil)
		wantAss := map[string][]*domain.CodesetAssignment{}
		if d := cmp.Diff(wantAss, gotAss); d != "" {
			t.Errorf("Unexpected Assignment: %s", diff.PrintWantGot(d))
		}
	})
}

func TestGetAllCodesetAssignments(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		want := make(map[string][]*domain.CodesetAssignment, len(codesets))

		addToWantAssignment := func(wf string, cs *domain.Codeset, webhookID *int64) {
			csAsg := domain.CodesetAssignment{Codeset: cs, WebhookID: webhookID}
			if a, exists := want[wf]; exists {
				want[wf] = append(a, &csAsg)
			} else {
				want[wf] = []*domain.CodesetAssignment{&csAsg}
			}
		}

		// wf0 -> not assigned
		// wf1 -> cs1
		// wf2 -> cs0, cs2
		for i := 0; i < len(codesets); i++ {
			wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: fmt.Sprintf("wf%d", i)})
			assertError(t, err, nil)
			if i != 0 {
				if i == 2 {
					cs := codesets[i-2]
					_, webhookID, err := mgr.AssignToCodeset(context.Background(), wf.Name, cs.Project, cs.Name)
					assertError(t, err, nil)
					addToWantAssignment(wf.Name, cs, webhookID)
				}
				_, webhookID, err := mgr.AssignToCodeset(context.Background(), wf.Name, codesets[i].Project, codesets[i].Name)
				assertError(t, err, nil)
				addToWantAssignment(wf.Name, codesets[i], webhookID)
			}
		}

		// list all
		got := mgr.GetAllCodesetAssignments(context.Background(), nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Assignments: %s", diff.PrintWantGot(d))
		}

		// list from a specific workflow
		for i := 0; i < len(codesets); i++ {
			wf := fmt.Sprintf("wf%d", i)
			got := mgr.GetAllCodesetAssignments(context.TODO(), &wf)

			if d := cmp.Diff(want[wf], got[wf]); d != "" {
				t.Errorf("Unexpected Workflow Assignments: %s", diff.PrintWantGot(d))
			}
		}
	})
}

func TestGetWorkflowRuns(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		want := []*domain.WorkflowRun{}

		// filter nil, no runs
		got, err := mgr.GetWorkflowRuns(context.Background(), nil)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// with filter, no runs
		filter := domain.WorkflowRunFilter{}
		got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		// create 3 runs with (cs0, csproject0, "Succeeded", "Failed", "Succeeded") and list
		for i := 0; i < 3; i++ {
			// currently, assigning a workflow to a codeset is the only function that creates a workflow run
			_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
			assertError(t, err, nil)

			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ = workflowBackend.GetWorkflowRuns(context.TODO(), wf, nil)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}
	})

	t.Run("filter by workflow", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		want := []*domain.WorkflowRun{}

		// non existing workflow, no runs
		wfName := "unknownWf"
		filterNoRunsNoWf := domain.WorkflowRunFilter{WorkflowName: &wfName}
		got, err := mgr.GetWorkflowRuns(context.Background(), &filterNoRunsNoWf)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// create a workflow
		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		// existing workflow, no runs
		filterNoRunsExistingWf := domain.WorkflowRunFilter{WorkflowName: &wf.Name}
		got, err = mgr.GetWorkflowRuns(context.Background(), &filterNoRunsExistingWf)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// create multiple workflows/runs:
		// wf0 -> 0 runs
		// wf1 -> 1 run (cs0, csproject0, Succeeded)
		// wf2 -> 2 runs (cs0, csproject0, Succeeded, Failed)
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		for i := 0; i < len(codesets); i++ {
			wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: fmt.Sprintf("wf%d", i)})
			assertError(t, err, nil)

			for j := 0; j < i; j++ {
				_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
				assertError(t, err, nil)
			}

		}

		// iterate over each workflow listing its runs
		workflows := workflowStore.GetWorkflows(context.TODO(), nil)
		for _, wf := range workflows {
			filter := domain.WorkflowRunFilter{WorkflowName: &wf.Name}
			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), &domain.Workflow{Name: wf.Name}, nil)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}
	})

	t.Run("filter by codeset", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		want := []*domain.WorkflowRun{}

		// non existing codeset, no runs
		csName := "unknownCs"
		filterNoRunsNoCs := domain.WorkflowRunFilter{CodesetName: csName}
		got, err := mgr.GetWorkflowRuns(context.Background(), &filterNoRunsNoCs)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// existing codeset, no runs
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		filterNoRuns := domain.WorkflowRunFilter{CodesetName: codesets[0].Name}
		got, err = mgr.GetWorkflowRuns(context.Background(), &filterNoRuns)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		// create multiple runs for the same workflow (wf) using different codesets:
		// 1. (cs0, csproject0, Succeeded)
		// 2. (cs1, csproject1, Failed)
		// 3. (cs2, csproject1, Succeeded)
		for i := 0; i < len(codesets); i++ {
			_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codesets[i].Project, codesets[i].Name)
			assertError(t, err, nil)
		}

		// iterate over each codeset and list runs by codeset name
		for _, cs := range codesets {
			filter := domain.WorkflowRunFilter{CodesetName: cs.Name}
			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), wf, &filter)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}

		// iterate over each codeset and list runs by codeset project
		for _, cs := range codesets {
			filter := domain.WorkflowRunFilter{CodesetProject: cs.Project}
			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), wf, &filter)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}

		// iterate over each codeset and list runs by codeset name and project
		for _, cs := range codesets {
			filter := domain.WorkflowRunFilter{CodesetName: cs.Name, CodesetProject: cs.Project}
			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), wf, &filter)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		want := []*domain.WorkflowRun{}

		// nil status, no runs
		filterNoRunsNilStatus := domain.WorkflowRunFilter{Status: nil}
		got, err := mgr.GetWorkflowRuns(context.Background(), &filterNoRunsNilStatus)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// empty status, no runs
		filterNoRunsEmptyStatus := domain.WorkflowRunFilter{Status: []string{}}
		got, err = mgr.GetWorkflowRuns(context.Background(), &filterNoRunsEmptyStatus)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		// with status, no runs
		filterNoRunsWithStatus := domain.WorkflowRunFilter{Status: []string{"Succeeded"}}
		got, err = mgr.GetWorkflowRuns(context.Background(), &filterNoRunsWithStatus)
		assertError(t, err, nil)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
		}

		wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		// create 3 runs for workflow 'wf' using cs0:
		// 1. (cs0, csproject0, Succeeded)
		// 2. (cs0, csproject0, Failed)
		// 3. (cs0, csproject0, Succeeded)
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		for i := 0; i < len(codesets); i++ {
			_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codesets[0].Project, codesets[0].Name)
			assertError(t, err, nil)
		}

		// iterate over the worklow statuses and list by it
		for _, s := range workflowRunStatuses {
			status := []string{s}
			filter := domain.WorkflowRunFilter{Status: status}
			got, err = mgr.GetWorkflowRuns(context.Background(), &filter)
			assertError(t, err, nil)

			want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), wf, &filter)
			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
			}
		}
	})

	t.Run("filter by workflow, codeset and status", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		// create multiple workflows, and runs using varying codesets
		// wf0 -> 0 runs
		// wf1 -> 1 run (cs0, project0, Succeeded)
		// wf2 -> 2 runs (cs1, project1, Succeeded) (cs2, project1, Failed)
		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		for i := 0; i < len(codesets); i++ {
			wf, err := mgr.CreateWorkflow(context.Background(), &domain.Workflow{Name: fmt.Sprintf("wf%d", i)})
			assertError(t, err, nil)

			for j := 0; j < i; j++ {
				csIndex := j
				if i == 2 {
					csIndex = j + 1
				}
				_, _, err = mgr.AssignToCodeset(context.Background(), wf.Name, codesets[csIndex].Project, codesets[csIndex].Name)
				assertError(t, err, nil)
			}
		}

		// iterate over all workflows, codesets, status listing runs and filtering for each combination
		workflows := workflowStore.GetWorkflows(context.TODO(), nil)
		for _, wf := range workflows {
			wfName := wf.Name
			for _, cs := range codesets {
				csName := cs.Name
				csProject := cs.Project
				for _, status := range workflowRunStatuses {
					status := []string{status}
					filter := domain.WorkflowRunFilter{WorkflowName: &wfName, CodesetName: csName, CodesetProject: csProject, Status: status}
					got, err := mgr.GetWorkflowRuns(context.Background(), &filter)
					assertError(t, err, nil)

					want, _ := workflowBackend.GetWorkflowRuns(context.TODO(), &domain.Workflow{Name: wfName}, &filter)
					if d := cmp.Diff(want, got); d != "" {
						t.Errorf("Unexpected Workflow Runs: %s", diff.PrintWantGot(d))
					}
				}
			}
		}
	})
}

func TestGetAssignmentStatus(t *testing.T) {
	t.Run("not assigned", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.TODO(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		got := mgr.GetAssignmentStatus(context.TODO(), wf.Name)
		want := domain.WorkflowAssignmentStatus{}
		if d := cmp.Diff(&want, got); d != "" {
			t.Errorf("Unexpected Workflow Assignment Status: %s", diff.PrintWantGot(d))
		}
	})

	t.Run("assigned", func(t *testing.T) {
		mgr := newFakeWorkflowManager(t)

		wf, err := mgr.CreateWorkflow(context.TODO(), &domain.Workflow{Name: "wf"})
		assertError(t, err, nil)

		codesets, _ := codesetStore.GetAll(context.TODO(), nil, nil)
		codeset := codesets[0]

		listener, _, err := mgr.AssignToCodeset(context.TODO(), wf.Name, codeset.Project, codeset.Name)
		assertError(t, err, nil)

		got := mgr.GetAssignmentStatus(context.TODO(), wf.Name)
		want := domain.WorkflowAssignmentStatus{Available: listener.Available, URL: listener.DashboardURL}
		if d := cmp.Diff(&want, got); d != "" {
			t.Errorf("Unexpected Workflow Assignment Status: %s", diff.PrintWantGot(d))
		}
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func newFakeWorkflowManager(t *testing.T) *WorkflowManager {
	t.Helper()

	workflowStore = core.NewWorkflowStore()
	workflowBackend = &fakeWorkflowBackend{t, make(map[string]*fakeStorableWorkflow)}
	codesetStore = &fakeCodesetStore{t, make(map[codesetID]fakeStorableCodeset)}
	extensionRegistry = NewExtensionRegistry(core.NewExtensionStore())

	// add codesets to the codeset store for the tests to use it:
	// 1. name: cs0, project: csproject0
	// 2. name: cs1, project: csproject1
	// 3. name: cs2, project: csproject1
	for i := 0; i < 3; i++ {
		projectIndex := i
		if i == 2 {
			projectIndex = i - 1
		}
		_, _, _, err := codesetStore.Add(context.Background(), &domain.Codeset{
			Name:    fmt.Sprintf("cs%d", i),
			Project: fmt.Sprintf("csproject%d", projectIndex),
			URL:     fmt.Sprintf("http://codeset/test-project%d/cs%d", projectIndex, i),
		})
		if err != nil {
			t.Fatalf("Error initializing fake codeset store")
		}
	}

	return NewWorkflowManager(workflowBackend, workflowStore, codesetStore, extensionRegistry)
}

func createFakeExtension(t *testing.T, wfm *WorkflowManager, prefix string) *domain.Extension {
	t.Helper()

	e := &domain.Extension{
		ID:      prefix + "extension",
		Product: prefix + "product",
		Version: "1.0",
		Zone:    prefix + "zone",
	}
	s1 := &domain.ExtensionService{
		ID:       prefix + "service-001",
		Resource: prefix + "resource",
		Category: prefix + "category",
	}
	e.AddService(s1)
	s2 := &domain.ExtensionService{
		ID: prefix + "service-002",
	}
	e.AddService(s2)
	ep1 := &domain.ExtensionServiceEndpoint{
		URL:  fmt.Sprintf("https://%sendpoint-001.com", prefix),
		Type: domain.EETInternal,
	}
	e.AddEndpoint(s1.ID, ep1)
	ep2 := &domain.ExtensionServiceEndpoint{
		URL:  fmt.Sprintf("https://%sendpoint-002.com", prefix),
		Type: domain.EETExternal,
	}
	e.AddEndpoint(s2.ID, ep2)
	c1 := &domain.ExtensionServiceCredentials{
		ID: prefix + "credentials-001",
	}
	e.AddCredentials(s1.ID, c1)
	c2 := &domain.ExtensionServiceCredentials{
		ID: prefix + "credentials-002",
	}
	e.AddCredentials(s2.ID, c2)

	return e
}

type fakeStorableWorkflow struct {
	listener *domain.WorkflowListener
	runs     []*domain.WorkflowRun
}

type fakeWorkflowBackend struct {
	t         *testing.T
	workflows map[string]*fakeStorableWorkflow
}

func (b *fakeWorkflowBackend) CreateWorkflow(ctx context.Context, w *domain.Workflow) error {
	b.t.Helper()

	if _, exists := b.workflows[w.Name]; exists {
		return domain.ErrWorkflowExists
	}
	b.workflows[w.Name] = &fakeStorableWorkflow{nil, []*domain.WorkflowRun{}}
	return nil
}

func (b *fakeWorkflowBackend) DeleteWorkflow(ctx context.Context, workflowName string) error {
	b.t.Helper()

	delete(b.workflows, workflowName)
	return nil
}

func (b *fakeWorkflowBackend) CreateWorkflowRun(ctx context.Context, workflowName string, codeset *domain.Codeset) error {
	b.t.Helper()

	if _, exists := b.workflows[workflowName]; !exists {
		return fmt.Errorf("workflow not found")
	}

	runs := b.workflows[workflowName].runs
	run := &domain.WorkflowRun{
		Name:        fmt.Sprintf("%s-run%d", workflowName, len(runs)),
		WorkflowRef: workflowName,
		Inputs: []*domain.WorkflowRunInput{
			{Input: &domain.WorkflowInput{Name: "codeset-name", Type: "codeset"}, Value: fmt.Sprintf("%s/%s", codeset.Project, codeset.Name)},
			{Input: &domain.WorkflowInput{Name: "predictor", Type: "string"}, Value: "sklearn"}},
		Status: workflowRunStatuses[len(runs)%len(workflowRunStatuses)]}

	b.workflows[workflowName].runs = append(b.workflows[workflowName].runs, run)
	return nil
}

func (b *fakeWorkflowBackend) GetWorkflowRuns(ctx context.Context, wf *domain.Workflow, filter *domain.WorkflowRunFilter) ([]*domain.WorkflowRun, error) {
	b.t.Helper()

	res := []*domain.WorkflowRun{}
	if sw, exists := b.workflows[wf.Name]; !exists || len(sw.runs) == 0 {
		return res, nil
	}

	runs := b.workflows[wf.Name].runs
	if filter == nil || (filter.CodesetName == "" && filter.CodesetProject == "" && len(filter.Status) == 0) {
		return runs, nil
	}

	getCodesetProjectName := func(inputValue string) (string, string) {
		nameAndValue := strings.Split(inputValue, "/")
		return nameAndValue[0], nameAndValue[1]
	}

	if filter.CodesetName != "" || filter.CodesetProject != "" {
		for _, run := range runs {
			if len(filter.Status) == 0 || util.StringInSlice(run.Status, filter.Status) {
				for _, input := range run.Inputs {
					if input.Input.Type == "codeset" {
						csProject, csName := getCodesetProjectName(input.Value)
						if filter.CodesetName != "" && filter.CodesetProject != "" {
							if filter.CodesetName == csName && filter.CodesetProject == csProject {
								res = append(res, run)
							}
						} else if (filter.CodesetName == csName && filter.CodesetProject == "") || (filter.CodesetProject == csProject && filter.CodesetName == "") {
							res = append(res, run)
						}
					}
				}
			}
		}
		return res, nil
	}

	if len(filter.Status) > 0 {
		for _, run := range runs {
			if util.StringInSlice(run.Status, filter.Status) {
				res = append(res, run)
			}
		}
	}

	return res, nil
}

func (b *fakeWorkflowBackend) CreateWorkflowListener(ctx context.Context, workflowName string, timeout time.Duration) (*domain.WorkflowListener, error) {
	b.t.Helper()

	listener := b.workflows[workflowName].listener
	if listener == nil {
		listener = &domain.WorkflowListener{Name: workflowName, Available: true, URL: fmt.Sprintf("http://%s.listener.test", workflowName),
			DashboardURL: fmt.Sprintf("http://dashboard.test/%s", workflowName)}
		b.workflows[workflowName].listener = listener
	}
	return listener, nil
}

func (b *fakeWorkflowBackend) DeleteWorkflowListener(ctx context.Context, workflowName string) error {
	b.t.Helper()

	wf, exists := b.workflows[workflowName]
	if !exists {
		return nil
	}

	wf.listener = nil
	return nil
}

func (b *fakeWorkflowBackend) GetWorkflowListener(ctx context.Context, workflowName string) (*domain.WorkflowListener, error) {
	b.t.Helper()

	if wf, exists := b.workflows[workflowName]; exists {
		if wf.listener != nil {
			return wf.listener, nil
		}
	}
	return nil, fmt.Errorf("listener not found")
}

type codesetID struct {
	name    string
	project string
}

type fakeStorableCodeset struct {
	codeset     *domain.Codeset
	webhooks    map[int64]string
	subscribers []domain.CodesetSubscriber
}

type fakeCodesetStore struct {
	t     *testing.T
	store map[codesetID]fakeStorableCodeset
}

func (fcs *fakeCodesetStore) Add(ctx context.Context, c *domain.Codeset) (*domain.Codeset, *string, *string, error) {
	fcs.t.Helper()

	fcs.store[codesetID{c.Name, c.Project}] = fakeStorableCodeset{codeset: c, webhooks: make(map[int64]string)}
	return c, nil, nil, nil
}

func (fcs *fakeCodesetStore) CreateWebhook(ctx context.Context, c *domain.Codeset, url string) (*int64, error) {
	fcs.t.Helper()

	id := rand.Int63()
	fcs.store[codesetID{c.Name, c.Project}].webhooks[id] = url
	return &id, nil
}

func (fcs *fakeCodesetStore) DeleteWebhook(ctx context.Context, c *domain.Codeset, id *int64) error {
	fcs.t.Helper()

	delete(fcs.store[codesetID{c.Name, c.Project}].webhooks, *id)
	return nil
}

func (fcs *fakeCodesetStore) Delete(ctx context.Context, project, name string) error {
	fcs.t.Helper()

	sc, ok := fcs.store[codesetID{name, project}]
	if !ok {
		return nil
	}
	for _, subscriber := range sc.subscribers {
		subscriber.OnDeletingCodeset(ctx, sc.codeset)
	}

	delete(fcs.store, codesetID{name, project})
	return nil
}

func (fcs *fakeCodesetStore) Find(ctx context.Context, project, name string) (*domain.Codeset, error) {
	fcs.t.Helper()

	if sc, exists := fcs.store[codesetID{name, project}]; exists {
		return sc.codeset, nil
	}
	return nil, errCodesetNotFound
}

func (fcs *fakeCodesetStore) GetAll(ctx context.Context, project, label *string) (res []*domain.Codeset, err error) {
	fcs.t.Helper()

	for _, c := range fcs.store {
		res = append(res, c.codeset)
	}
	return res, nil
}

func (fcs *fakeCodesetStore) Subscribe(ctx context.Context, subscriber domain.CodesetSubscriber, codeset *domain.Codeset) error {
	fcs.t.Helper()

	sc, ok := fcs.store[codesetID{codeset.Name, codeset.Project}]
	if !ok {
		return fmt.Errorf("codeset not found")
	}
	sc.subscribers = append(sc.subscribers, subscriber)
	fcs.store[codesetID{codeset.Name, codeset.Project}] = sc
	return nil
}

func (fcs *fakeCodesetStore) Unsubscribe(ctx context.Context, subscriber domain.CodesetSubscriber, codeset *domain.Codeset) error {
	fcs.t.Helper()

	sc := fcs.store[codesetID{codeset.Name, codeset.Project}]
	sc.subscribers = removeSubscriber(sc.subscribers, subscriber)
	fcs.store[codesetID{codeset.Name, codeset.Project}] = sc
	return nil
}

func (fcs *fakeCodesetStore) getSubscribers(ctx context.Context, c *domain.Codeset) []domain.CodesetSubscriber {
	return fcs.store[codesetID{c.Name, c.Project}].subscribers
}

func removeSubscriber(subscribers []domain.CodesetSubscriber, subscriber domain.CodesetSubscriber) []domain.CodesetSubscriber {
	for i, s := range subscribers {
		if s == subscriber {
			subscribers[len(subscribers)-1], subscribers[i] = subscribers[i], subscribers[len(subscribers)-1]
			return subscribers[:len(subscribers)-1]
		}
	}
	return subscribers
}
