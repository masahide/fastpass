package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// GetTicketHandlerFunc turns a function with the right signature into a get ticket handler
type GetTicketHandlerFunc func(GetTicketParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTicketHandlerFunc) Handle(params GetTicketParams) middleware.Responder {
	return fn(params)
}

// GetTicketHandler interface for that can handle valid get ticket params
type GetTicketHandler interface {
	Handle(GetTicketParams) middleware.Responder
}

// NewGetTicket creates a new http.Handler for the get ticket operation
func NewGetTicket(ctx *middleware.Context, handler GetTicketHandler) *GetTicket {
	return &GetTicket{Context: ctx, Handler: handler}
}

/*GetTicket swagger:route GET /events/{id}/tickets/{uid} getTicket

チケット確認

*/
type GetTicket struct {
	Context *middleware.Context
	Params  GetTicketParams
	Handler GetTicketHandler
}

func (o *GetTicket) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	o.Params = NewGetTicketParams()

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
