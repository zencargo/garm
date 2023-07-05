// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new repositories API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for repositories API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateRepo(params *CreateRepoParams, opts ...ClientOption) (*CreateRepoOK, error)

	CreateRepoPool(params *CreateRepoPoolParams, opts ...ClientOption) (*CreateRepoPoolOK, error)

	DeleteRepo(params *DeleteRepoParams, opts ...ClientOption) error

	DeleteRepoPool(params *DeleteRepoPoolParams, opts ...ClientOption) error

	GetRepo(params *GetRepoParams, opts ...ClientOption) (*GetRepoOK, error)

	GetRepoPool(params *GetRepoPoolParams, opts ...ClientOption) (*GetRepoPoolOK, error)

	ListRepoInstances(params *ListRepoInstancesParams, opts ...ClientOption) (*ListRepoInstancesOK, error)

	ListRepoPools(params *ListRepoPoolsParams, opts ...ClientOption) (*ListRepoPoolsOK, error)

	ListRepos(params *ListReposParams, opts ...ClientOption) (*ListReposOK, error)

	UpdateRepo(params *UpdateRepoParams, opts ...ClientOption) (*UpdateRepoOK, error)

	UpdateRepoPool(params *UpdateRepoPoolParams, opts ...ClientOption) (*UpdateRepoPoolOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
CreateRepo creates repository with the parameters given
*/
func (a *Client) CreateRepo(params *CreateRepoParams, opts ...ClientOption) (*CreateRepoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateRepoParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "CreateRepo",
		Method:             "POST",
		PathPattern:        "/repositories",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateRepoReader{formats: a.formats},
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
	success, ok := result.(*CreateRepoOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateRepoDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
CreateRepoPool creates repository pool with the parameters given
*/
func (a *Client) CreateRepoPool(params *CreateRepoPoolParams, opts ...ClientOption) (*CreateRepoPoolOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateRepoPoolParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "CreateRepoPool",
		Method:             "POST",
		PathPattern:        "/repositories/{repoID}/pools",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateRepoPoolReader{formats: a.formats},
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
	success, ok := result.(*CreateRepoPoolOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateRepoPoolDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
DeleteRepo deletes repository by ID
*/
func (a *Client) DeleteRepo(params *DeleteRepoParams, opts ...ClientOption) error {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteRepoParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteRepo",
		Method:             "DELETE",
		PathPattern:        "/repositories/{repoID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteRepoReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	_, err := a.transport.Submit(op)
	if err != nil {
		return err
	}
	return nil
}

/*
DeleteRepoPool deletes repository pool by ID
*/
func (a *Client) DeleteRepoPool(params *DeleteRepoPoolParams, opts ...ClientOption) error {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteRepoPoolParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteRepoPool",
		Method:             "DELETE",
		PathPattern:        "/repositories/{repoID}/pools/{poolID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteRepoPoolReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	_, err := a.transport.Submit(op)
	if err != nil {
		return err
	}
	return nil
}

/*
GetRepo gets repository by ID
*/
func (a *Client) GetRepo(params *GetRepoParams, opts ...ClientOption) (*GetRepoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRepoParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRepo",
		Method:             "GET",
		PathPattern:        "/repositories/{repoID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetRepoReader{formats: a.formats},
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
	success, ok := result.(*GetRepoOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetRepoDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetRepoPool gets repository pool by ID
*/
func (a *Client) GetRepoPool(params *GetRepoPoolParams, opts ...ClientOption) (*GetRepoPoolOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRepoPoolParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRepoPool",
		Method:             "GET",
		PathPattern:        "/repositories/{repoID}/pools/{poolID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetRepoPoolReader{formats: a.formats},
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
	success, ok := result.(*GetRepoPoolOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetRepoPoolDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListRepoInstances lists repository instances
*/
func (a *Client) ListRepoInstances(params *ListRepoInstancesParams, opts ...ClientOption) (*ListRepoInstancesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListRepoInstancesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListRepoInstances",
		Method:             "GET",
		PathPattern:        "/repositories/{repoID}/instances",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListRepoInstancesReader{formats: a.formats},
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
	success, ok := result.(*ListRepoInstancesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListRepoInstancesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListRepoPools lists repository pools
*/
func (a *Client) ListRepoPools(params *ListRepoPoolsParams, opts ...ClientOption) (*ListRepoPoolsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListRepoPoolsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListRepoPools",
		Method:             "GET",
		PathPattern:        "/repositories/{repoID}/pools",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListRepoPoolsReader{formats: a.formats},
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
	success, ok := result.(*ListRepoPoolsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListRepoPoolsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListRepos lists repositories
*/
func (a *Client) ListRepos(params *ListReposParams, opts ...ClientOption) (*ListReposOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListReposParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListRepos",
		Method:             "GET",
		PathPattern:        "/repositories",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListReposReader{formats: a.formats},
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
	success, ok := result.(*ListReposOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListReposDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
UpdateRepo updates repository with the parameters given
*/
func (a *Client) UpdateRepo(params *UpdateRepoParams, opts ...ClientOption) (*UpdateRepoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateRepoParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "UpdateRepo",
		Method:             "PUT",
		PathPattern:        "/repositories/{repoID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateRepoReader{formats: a.formats},
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
	success, ok := result.(*UpdateRepoOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*UpdateRepoDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
UpdateRepoPool updates repository pool with the parameters given
*/
func (a *Client) UpdateRepoPool(params *UpdateRepoPoolParams, opts ...ClientOption) (*UpdateRepoPoolOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateRepoPoolParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "UpdateRepoPool",
		Method:             "PUT",
		PathPattern:        "/repositories/{repoID}/pools/{poolID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateRepoPoolReader{formats: a.formats},
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
	success, ok := result.(*UpdateRepoPoolOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*UpdateRepoPoolDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
