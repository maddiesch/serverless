package middleware

import (
	"context"
	"net/http"
)

type Aborter func(error)

type Handler func(context.Context, http.ResponseWriter, *http.Request, Aborter)
