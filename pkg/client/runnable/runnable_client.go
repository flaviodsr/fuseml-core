// Code generated by go-swagger; DO NOT EDIT.

package runnable

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new runnable API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for runnable API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetRunnable(params *GetRunnableParams, opts ...ClientOption) (*GetRunnableOK, error)

	ListRunnables(params *ListRunnablesParams, opts ...ClientOption) (*ListRunnablesOK, error)

	RegisterRunnable(params *RegisterRunnableParams, opts ...ClientOption) (*RegisterRunnableCreated, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetRunnable gets information about a runnable

  Lookup information about a runnable in the FuseML store
*/
func (a *Client) GetRunnable(params *GetRunnableParams, opts ...ClientOption) (*GetRunnableOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRunnableParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getRunnable",
		Method:             "GET",
		PathPattern:        "/runnables/{runnableNameOrId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetRunnableReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetRunnableOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getRunnable: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  ListRunnables retrieves runnable information

  Get information about runnables currently registered in the FuseML runnable store
*/
func (a *Client) ListRunnables(params *ListRunnablesParams, opts ...ClientOption) (*ListRunnablesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListRunnablesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listRunnables",
		Method:             "GET",
		PathPattern:        "/runnables",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListRunnablesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListRunnablesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for listRunnables: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  RegisterRunnable registers runnable

  Register a runnable with the FuseML runnable store
*/
func (a *Client) RegisterRunnable(params *RegisterRunnableParams, opts ...ClientOption) (*RegisterRunnableCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRegisterRunnableParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "registerRunnable",
		Method:             "POST",
		PathPattern:        "/runnables",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &RegisterRunnableReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*RegisterRunnableCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for registerRunnable: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
