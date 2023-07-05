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

// ListRepoInstancesReader is a Reader for the ListRepoInstances structure.
type ListRepoInstancesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListRepoInstancesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListRepoInstancesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListRepoInstancesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListRepoInstancesOK creates a ListRepoInstancesOK with default headers values
func NewListRepoInstancesOK() *ListRepoInstancesOK {
	return &ListRepoInstancesOK{}
}

/*
ListRepoInstancesOK describes a response with status code 200, with default header values.

Instances
*/
type ListRepoInstancesOK struct {
	Payload garm_params.Instances
}

// IsSuccess returns true when this list repo instances o k response has a 2xx status code
func (o *ListRepoInstancesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list repo instances o k response has a 3xx status code
func (o *ListRepoInstancesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list repo instances o k response has a 4xx status code
func (o *ListRepoInstancesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list repo instances o k response has a 5xx status code
func (o *ListRepoInstancesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list repo instances o k response a status code equal to that given
func (o *ListRepoInstancesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list repo instances o k response
func (o *ListRepoInstancesOK) Code() int {
	return 200
}

func (o *ListRepoInstancesOK) Error() string {
	return fmt.Sprintf("[GET /repositories/{repoID}/instances][%d] listRepoInstancesOK  %+v", 200, o.Payload)
}

func (o *ListRepoInstancesOK) String() string {
	return fmt.Sprintf("[GET /repositories/{repoID}/instances][%d] listRepoInstancesOK  %+v", 200, o.Payload)
}

func (o *ListRepoInstancesOK) GetPayload() garm_params.Instances {
	return o.Payload
}

func (o *ListRepoInstancesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRepoInstancesDefault creates a ListRepoInstancesDefault with default headers values
func NewListRepoInstancesDefault(code int) *ListRepoInstancesDefault {
	return &ListRepoInstancesDefault{
		_statusCode: code,
	}
}

/*
ListRepoInstancesDefault describes a response with status code -1, with default header values.

APIErrorResponse
*/
type ListRepoInstancesDefault struct {
	_statusCode int

	Payload apiserver_params.APIErrorResponse
}

// IsSuccess returns true when this list repo instances default response has a 2xx status code
func (o *ListRepoInstancesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list repo instances default response has a 3xx status code
func (o *ListRepoInstancesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list repo instances default response has a 4xx status code
func (o *ListRepoInstancesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list repo instances default response has a 5xx status code
func (o *ListRepoInstancesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list repo instances default response a status code equal to that given
func (o *ListRepoInstancesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list repo instances default response
func (o *ListRepoInstancesDefault) Code() int {
	return o._statusCode
}

func (o *ListRepoInstancesDefault) Error() string {
	return fmt.Sprintf("[GET /repositories/{repoID}/instances][%d] ListRepoInstances default  %+v", o._statusCode, o.Payload)
}

func (o *ListRepoInstancesDefault) String() string {
	return fmt.Sprintf("[GET /repositories/{repoID}/instances][%d] ListRepoInstances default  %+v", o._statusCode, o.Payload)
}

func (o *ListRepoInstancesDefault) GetPayload() apiserver_params.APIErrorResponse {
	return o.Payload
}

func (o *ListRepoInstancesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
