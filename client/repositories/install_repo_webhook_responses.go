// Code generated by go-swagger; DO NOT EDIT.

package repositories

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	apiserver_params "github.com/cloudbase/garm/apiserver/params"
	garm_params "github.com/cloudbase/garm/params"
)

// InstallRepoWebhookReader is a Reader for the InstallRepoWebhook structure.
type InstallRepoWebhookReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InstallRepoWebhookReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewInstallRepoWebhookOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewInstallRepoWebhookDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewInstallRepoWebhookOK creates a InstallRepoWebhookOK with default headers values
func NewInstallRepoWebhookOK() *InstallRepoWebhookOK {
	return &InstallRepoWebhookOK{}
}

/*
InstallRepoWebhookOK describes a response with status code 200, with default header values.

HookInfo
*/
type InstallRepoWebhookOK struct {
	Payload garm_params.HookInfo
}

// IsSuccess returns true when this install repo webhook o k response has a 2xx status code
func (o *InstallRepoWebhookOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this install repo webhook o k response has a 3xx status code
func (o *InstallRepoWebhookOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this install repo webhook o k response has a 4xx status code
func (o *InstallRepoWebhookOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this install repo webhook o k response has a 5xx status code
func (o *InstallRepoWebhookOK) IsServerError() bool {
	return false
}

// IsCode returns true when this install repo webhook o k response a status code equal to that given
func (o *InstallRepoWebhookOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the install repo webhook o k response
func (o *InstallRepoWebhookOK) Code() int {
	return 200
}

func (o *InstallRepoWebhookOK) Error() string {
	return fmt.Sprintf("[POST /repositories/{repoID}/webhook][%d] installRepoWebhookOK  %+v", 200, o.Payload)
}

func (o *InstallRepoWebhookOK) String() string {
	return fmt.Sprintf("[POST /repositories/{repoID}/webhook][%d] installRepoWebhookOK  %+v", 200, o.Payload)
}

func (o *InstallRepoWebhookOK) GetPayload() garm_params.HookInfo {
	return o.Payload
}

func (o *InstallRepoWebhookOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstallRepoWebhookDefault creates a InstallRepoWebhookDefault with default headers values
func NewInstallRepoWebhookDefault(code int) *InstallRepoWebhookDefault {
	return &InstallRepoWebhookDefault{
		_statusCode: code,
	}
}

/*
InstallRepoWebhookDefault describes a response with status code -1, with default header values.

APIErrorResponse
*/
type InstallRepoWebhookDefault struct {
	_statusCode int

	Payload apiserver_params.APIErrorResponse
}

// IsSuccess returns true when this install repo webhook default response has a 2xx status code
func (o *InstallRepoWebhookDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this install repo webhook default response has a 3xx status code
func (o *InstallRepoWebhookDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this install repo webhook default response has a 4xx status code
func (o *InstallRepoWebhookDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this install repo webhook default response has a 5xx status code
func (o *InstallRepoWebhookDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this install repo webhook default response a status code equal to that given
func (o *InstallRepoWebhookDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the install repo webhook default response
func (o *InstallRepoWebhookDefault) Code() int {
	return o._statusCode
}

func (o *InstallRepoWebhookDefault) Error() string {
	return fmt.Sprintf("[POST /repositories/{repoID}/webhook][%d] InstallRepoWebhook default  %+v", o._statusCode, o.Payload)
}

func (o *InstallRepoWebhookDefault) String() string {
	return fmt.Sprintf("[POST /repositories/{repoID}/webhook][%d] InstallRepoWebhook default  %+v", o._statusCode, o.Payload)
}

func (o *InstallRepoWebhookDefault) GetPayload() apiserver_params.APIErrorResponse {
	return o.Payload
}

func (o *InstallRepoWebhookDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
