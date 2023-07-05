// Code generated by go-swagger; DO NOT EDIT.

package instances

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
)

// NewGetInstanceParams creates a new GetInstanceParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetInstanceParams() *GetInstanceParams {
	return &GetInstanceParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetInstanceParamsWithTimeout creates a new GetInstanceParams object
// with the ability to set a timeout on a request.
func NewGetInstanceParamsWithTimeout(timeout time.Duration) *GetInstanceParams {
	return &GetInstanceParams{
		timeout: timeout,
	}
}

// NewGetInstanceParamsWithContext creates a new GetInstanceParams object
// with the ability to set a context for a request.
func NewGetInstanceParamsWithContext(ctx context.Context) *GetInstanceParams {
	return &GetInstanceParams{
		Context: ctx,
	}
}

// NewGetInstanceParamsWithHTTPClient creates a new GetInstanceParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetInstanceParamsWithHTTPClient(client *http.Client) *GetInstanceParams {
	return &GetInstanceParams{
		HTTPClient: client,
	}
}

/*
GetInstanceParams contains all the parameters to send to the API endpoint

	for the get instance operation.

	Typically these are written to a http.Request.
*/
type GetInstanceParams struct {

	/* InstanceName.

	   Runner instance name.
	*/
	InstanceName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get instance params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetInstanceParams) WithDefaults() *GetInstanceParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get instance params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetInstanceParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get instance params
func (o *GetInstanceParams) WithTimeout(timeout time.Duration) *GetInstanceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get instance params
func (o *GetInstanceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get instance params
func (o *GetInstanceParams) WithContext(ctx context.Context) *GetInstanceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get instance params
func (o *GetInstanceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get instance params
func (o *GetInstanceParams) WithHTTPClient(client *http.Client) *GetInstanceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get instance params
func (o *GetInstanceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithInstanceName adds the instanceName to the get instance params
func (o *GetInstanceParams) WithInstanceName(instanceName string) *GetInstanceParams {
	o.SetInstanceName(instanceName)
	return o
}

// SetInstanceName adds the instanceName to the get instance params
func (o *GetInstanceParams) SetInstanceName(instanceName string) {
	o.InstanceName = instanceName
}

// WriteToRequest writes these params to a swagger request
func (o *GetInstanceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param instanceName
	if err := r.SetPathParam("instanceName", o.InstanceName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
