// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	garm_params "github.com/cloudbase/garm/params"
)

// NewUpdateRepoParams creates a new UpdateRepoParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateRepoParams() *UpdateRepoParams {
	return &UpdateRepoParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateRepoParamsWithTimeout creates a new UpdateRepoParams object
// with the ability to set a timeout on a request.
func NewUpdateRepoParamsWithTimeout(timeout time.Duration) *UpdateRepoParams {
	return &UpdateRepoParams{
		timeout: timeout,
	}
}

// NewUpdateRepoParamsWithContext creates a new UpdateRepoParams object
// with the ability to set a context for a request.
func NewUpdateRepoParamsWithContext(ctx context.Context) *UpdateRepoParams {
	return &UpdateRepoParams{
		Context: ctx,
	}
}

// NewUpdateRepoParamsWithHTTPClient creates a new UpdateRepoParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateRepoParamsWithHTTPClient(client *http.Client) *UpdateRepoParams {
	return &UpdateRepoParams{
		HTTPClient: client,
	}
}

/*
UpdateRepoParams contains all the parameters to send to the API endpoint

	for the update repo operation.

	Typically these are written to a http.Request.
*/
type UpdateRepoParams struct {

	/* Body.

	   Parameters used when updating the repository.
	*/
	Body garm_params.UpdateEntityParams

	/* RepoID.

	   ID of the repository to update.
	*/
	RepoID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update repo params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateRepoParams) WithDefaults() *UpdateRepoParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update repo params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateRepoParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update repo params
func (o *UpdateRepoParams) WithTimeout(timeout time.Duration) *UpdateRepoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update repo params
func (o *UpdateRepoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update repo params
func (o *UpdateRepoParams) WithContext(ctx context.Context) *UpdateRepoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update repo params
func (o *UpdateRepoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update repo params
func (o *UpdateRepoParams) WithHTTPClient(client *http.Client) *UpdateRepoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update repo params
func (o *UpdateRepoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update repo params
func (o *UpdateRepoParams) WithBody(body garm_params.UpdateEntityParams) *UpdateRepoParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update repo params
func (o *UpdateRepoParams) SetBody(body garm_params.UpdateEntityParams) {
	o.Body = body
}

// WithRepoID adds the repoID to the update repo params
func (o *UpdateRepoParams) WithRepoID(repoID string) *UpdateRepoParams {
	o.SetRepoID(repoID)
	return o
}

// SetRepoID adds the repoId to the update repo params
func (o *UpdateRepoParams) SetRepoID(repoID string) {
	o.RepoID = repoID
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateRepoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param repoID
	if err := r.SetPathParam("repoID", o.RepoID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
