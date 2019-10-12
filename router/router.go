package router

import (
	"context"
	"net/http"

	"github.com/maddiesch/serverless/middleware"
)

// HandlerFunc is the function that will be called when a request matches the specific route.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// Handler defines the methods needed to create a handler
type Handler interface {
	// Handle adds a handler to a path.
	// Method, Path, Handler
	Handle(string, string, HandlerFunc)

	Group(string, ...middleware.Handler) Group
}

// Router is the interface that must be implemented to create a router
type Router interface {
	Handler

	Dispatch(context.Context, http.ResponseWriter, *http.Request) error
}

// Group contains a router handler
type Group interface {
	Handler

	Parent() Group

	Router() Router

	Middleware() []middleware.Handler
}
