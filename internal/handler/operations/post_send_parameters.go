// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewPostSendParams creates a new PostSendParams object
// no default values defined in spec.
func NewPostSendParams() PostSendParams {

	return PostSendParams{}
}

// PostSendParams contains all the bound params for the post send operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostSend
type PostSendParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*This message will be send all subscribers
	  Required: true
	  In: query
	*/
	Msg string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostSendParams() beforehand.
func (o *PostSendParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qMsg, qhkMsg, _ := qs.GetOK("msg")
	if err := o.bindMsg(qMsg, qhkMsg, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindMsg binds and validates parameter Msg from query.
func (o *PostSendParams) bindMsg(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("msg", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("msg", "query", raw); err != nil {
		return err
	}

	o.Msg = raw

	return nil
}
