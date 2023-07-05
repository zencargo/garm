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

// ListReposReader is a Reader for the ListRepos structure.
type ListReposReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListReposReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListReposOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListReposDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListReposOK creates a ListReposOK with default headers values
func NewListReposOK() *ListReposOK {
	return &ListReposOK{}
}

/*
ListReposOK describes a response with status code 200, with default header values.

Repositories
*/
type ListReposOK struct {
	Payload garm_params.Repositories
}

// IsSuccess returns true when this list repos o k response has a 2xx status code
func (o *ListReposOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list repos o k response has a 3xx status code
func (o *ListReposOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list repos o k response has a 4xx status code
func (o *ListReposOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list repos o k response has a 5xx status code
func (o *ListReposOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list repos o k response a status code equal to that given
func (o *ListReposOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list repos o k response
func (o *ListReposOK) Code() int {
	return 200
}

func (o *ListReposOK) Error() string {
	return fmt.Sprintf("[GET /repositories][%d] listReposOK  %+v", 200, o.Payload)
}

func (o *ListReposOK) String() string {
	return fmt.Sprintf("[GET /repositories][%d] listReposOK  %+v", 200, o.Payload)
}

func (o *ListReposOK) GetPayload() garm_params.Repositories {
	return o.Payload
}

func (o *ListReposOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListReposDefault creates a ListReposDefault with default headers values
func NewListReposDefault(code int) *ListReposDefault {
	return &ListReposDefault{
		_statusCode: code,
	}
}

/*
ListReposDefault describes a response with status code -1, with default header values.

APIErrorResponse
*/
type ListReposDefault struct {
	_statusCode int

	Payload apiserver_params.APIErrorResponse
}

// IsSuccess returns true when this list repos default response has a 2xx status code
func (o *ListReposDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list repos default response has a 3xx status code
func (o *ListReposDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list repos default response has a 4xx status code
func (o *ListReposDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list repos default response has a 5xx status code
func (o *ListReposDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list repos default response a status code equal to that given
func (o *ListReposDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list repos default response
func (o *ListReposDefault) Code() int {
	return o._statusCode
}

func (o *ListReposDefault) Error() string {
	return fmt.Sprintf("[GET /repositories][%d] ListRepos default  %+v", o._statusCode, o.Payload)
}

func (o *ListReposDefault) String() string {
	return fmt.Sprintf("[GET /repositories][%d] ListRepos default  %+v", o._statusCode, o.Payload)
}

func (o *ListReposDefault) GetPayload() apiserver_params.APIErrorResponse {
	return o.Payload
}

func (o *ListReposDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
