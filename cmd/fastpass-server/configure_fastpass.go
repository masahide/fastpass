package main

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/masahide/fastpass"
	"github.com/masahide/fastpass/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.FastpassAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	f := fastpass.NewFastPass()

	api.AddEventHandler = operations.AddEventHandlerFunc(f.AddEventHandler)
	api.DeleteEventHandler = operations.DeleteEventHandlerFunc(f.DeleteEventHandler)
	api.GetTicketHandler = operations.GetTicketHandlerFunc(f.GetTicketHandler)
	api.ListEventsHandler = operations.ListEventsHandlerFunc(f.ListEventsHandler)
	api.TicketingHandler = operations.TicketingHandlerFunc(f.TicketingHandler)

	api.ServerShutdown = func() {}
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
