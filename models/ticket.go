package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/validate"
	"github.com/go-swagger/go-swagger/strfmt"
)

/*Ticket Ticket ticket

swagger:model Ticket
*/
type Ticket struct {

	/* EventID event id

	Required: true
	*/
	EventID string `json:"event_id,omitempty"`

	/* Num num

	Required: true
	*/
	Num int32 `json:"num,omitempty"`
}

// Validate validates this ticket
func (m *Ticket) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEventID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateNum(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Ticket) validateEventID(formats strfmt.Registry) error {

	if err := validate.RequiredString("event_id", "body", string(m.EventID)); err != nil {
		return err
	}

	return nil
}

func (m *Ticket) validateNum(formats strfmt.Registry) error {

	if err := validate.Required("num", "body", int32(m.Num)); err != nil {
		return err
	}

	return nil
}
