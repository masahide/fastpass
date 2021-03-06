package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/masahide/fastpass/models"
)

/*AddEventOK 作成成功

swagger:response addEventOK
*/
type AddEventOK struct {

	// In: body
	Payload *models.Event `json:"body,omitempty"`
}

// NewAddEventOK creates AddEventOK with default headers values
func NewAddEventOK() *AddEventOK {
	return &AddEventOK{}
}

// WithPayload adds the payload to the add event o k response
func (o *AddEventOK) WithPayload(payload *models.Event) *AddEventOK {
	o.Payload = payload
	return o
}

// WriteResponse to the client
func (o *AddEventOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddEventConflict Conflict(作成済み)

swagger:response addEventConflict
*/
type AddEventConflict struct {
}

// NewAddEventConflict creates AddEventConflict with default headers values
func NewAddEventConflict() *AddEventConflict {
	return &AddEventConflict{}
}

// WriteResponse to the client
func (o *AddEventConflict) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(409)
}

/*AddEventInternalServerError unexpected error

swagger:response addEventInternalServerError
*/
type AddEventInternalServerError struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddEventInternalServerError creates AddEventInternalServerError with default headers values
func NewAddEventInternalServerError() *AddEventInternalServerError {
	return &AddEventInternalServerError{}
}

// WithPayload adds the payload to the add event internal server error response
func (o *AddEventInternalServerError) WithPayload(payload *models.Error) *AddEventInternalServerError {
	o.Payload = payload
	return o
}

// WriteResponse to the client
func (o *AddEventInternalServerError) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
