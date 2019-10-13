package middleware

import (
	"context"
	"net/http"
)

// Aborter is the function passed into the middleware handler that if called
// will immediately abort the request processing chain.
type Aborter func(error)

// Handler is a middleware handler
type Handler func(context.Context, http.ResponseWriter, *http.Request, Aborter) context.Context
