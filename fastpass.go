package fastpass

import (
	"sync"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/masahide/fastpass/models"
	"github.com/masahide/fastpass/restapi/operations"
)

type Fastpass struct {
	mu     sync.RWMutex
	Events map[string]Event
}

func NewFastPass() *Fastpass {
	return &Fastpass{Events: map[string]Event{}}
}

type Event struct {
	models.Event

	mu      sync.RWMutex
	Tickets map[string]*Ticket
	tickets []Ticket
	Current int32
}

type Ticket struct {
	models.Ticket
}

func (f *Fastpass) AddEventHandler(params operations.AddEventParams) middleware.Responder {
	f.mu.RLock()
	_, ok := f.Events[params.ID]
	f.mu.RUnlock()
	if ok {
		return operations.NewAddEventConflict()
	}
	e := Event{
		Event:   models.Event{ID: params.ID, MaxTicket: params.Options.MaxTicket},
		Tickets: map[string]*Ticket{},
		tickets: make([]Ticket, params.Options.MaxTicket),
	}
	f.mu.Lock()
	f.Events[params.ID] = e
	f.mu.Unlock()
	return operations.NewAddEventOK()
}
func (f *Fastpass) DeleteEventHandler(params operations.DeleteEventParams) middleware.Responder {
	f.mu.RLock()
	_, ok := f.Events[params.ID]
	f.mu.RUnlock()
	if !ok {
		return operations.NewDeleteEventInternalServerError().WithPayload(&models.Error{Code: 1, Message: "event not found."})
	}
	return operations.NewDeleteEventNoContent()
}

func (f *Fastpass) GetTicketHandler(params operations.GetTicketParams) middleware.Responder {
	f.mu.RLock()
	e, ok := f.Events[params.ID]
	f.mu.RUnlock()
	if !ok {
		return operations.NewGetTicketNotFound()
	}
	e.mu.RLock()
	ticket, ok := e.Tickets[params.UID]
	e.mu.RUnlock()
	if !ok {
		return operations.NewGetTicketNotFound()
	}
	return operations.NewGetTicketOK().WithPayload(&ticket.Ticket)
}
func (f *Fastpass) ListEventsHandler(params operations.ListEventsParams) middleware.Responder {
	f.mu.RLock()

	events := make([]*models.Event, len(f.Events))
	i := 0
	for _, e := range f.Events {
		events[i] = &e.Event
	}
	f.mu.RUnlock()
	return operations.NewListEventsOK().WithPayload(events)
}
func (f *Fastpass) TicketingHandler(params operations.TicketingParams) middleware.Responder {
	f.mu.RLock()
	e, ok := f.Events[params.ID]
	f.mu.RUnlock()
	if !ok {
		return operations.NewTicketingNotFound()
	}
	e.mu.RLock()
	_, ok = e.Tickets[params.UID]
	e.mu.RUnlock()
	if ok {
		return operations.NewTicketingConflict()
	}
	e.mu.Lock()
	num := e.Current
	t := &e.tickets[num]
	e.Tickets[params.UID] = t
	e.mu.Unlock()
	t.Ticket.EventID = params.ID
	t.Ticket.Num = num
	return operations.NewTicketingOK().WithPayload(&t.Ticket)
}
