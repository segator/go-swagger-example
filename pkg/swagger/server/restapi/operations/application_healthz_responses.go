// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ApplicationHealthzOKCode is the HTTP code returned for type ApplicationHealthzOK
const ApplicationHealthzOKCode int = 200

/*ApplicationHealthzOK OK message.

swagger:response applicationHealthzOK
*/
type ApplicationHealthzOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewApplicationHealthzOK creates ApplicationHealthzOK with default headers values
func NewApplicationHealthzOK() *ApplicationHealthzOK {

	return &ApplicationHealthzOK{}
}

// WithPayload adds the payload to the application healthz o k response
func (o *ApplicationHealthzOK) WithPayload(payload string) *ApplicationHealthzOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the application healthz o k response
func (o *ApplicationHealthzOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ApplicationHealthzOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
