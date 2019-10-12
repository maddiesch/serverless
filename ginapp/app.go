package ginapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maddiesch/serverless"
	"github.com/maddiesch/serverless/middleware"
	"github.com/maddiesch/serverless/proxy"
	"github.com/maddiesch/serverless/responder"
	"github.com/maddiesch/serverless/router"
	"github.com/maddiesch/serverless/sam"
)

func init() {
	if !sam.IsLocal() {
		gin.SetMode(gin.ReleaseMode)
	}
}

// New returns a new gin
func New() *App {
	engine := gin.New()

	engine.Use(contextAssigner)
	engine.Use(gin.Recovery())
	engine.Use(middleware.Recovery(recoveryHandler))
	engine.Use(middleware.Runtime())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger())

	engine.NoRoute(defaultNoRouteHandler)

	return &App{
		engine: engine,
		proxy:  &ginProxy{engine: engine},
		router: &ginRouter{engine: engine},
	}
}

const (
	ginSubContextKey = ".gin-sub-context"
)

// App is the default gin application for serverless
type App struct {
	serverless.Application

	engine *gin.Engine
	proxy  *ginProxy
	router *ginRouter
}

// Middleware returns the middleware for processing the incoming requests
func (a *App) Middleware() []middleware.Handler {
	return []middleware.Handler{}
}

// Adapter returns the lambda proxy
func (a *App) Adapter() proxy.Lambda {
	return a.proxy
}

// Router returns the default router from the config
func (a *App) Router() router.Router {
	return a.router
}

// ConfigureGin performs the function passed in.
// The passed function allows access to the underlying gin engine.
func (a *App) ConfigureGin(fn func(*gin.Engine)) {
	fn(a.engine)
}

func recoveryHandler(c *gin.Context, e interface{}) error {
	var details string

	switch value := e.(type) {
	case error:
		details = fmt.Sprintf("%v", value)
	case string:
		details = value
	default:
		panic(e)
	}

	meta := make(map[string]interface{})
	if serverless.IsDebug() {
		meta["ErrorDetails"] = fmt.Sprintf("%v", details)
	}

	return serverless.ErrorResponse(c, []*serverless.Error{
		&serverless.Error{
			Status:      "500",
			Code:        "unknown_error",
			Title:       "Unknown Error",
			Description: "An unknown error occurred. Please try your request again",
			Meta:        &meta,
		},
	})
}

func contextAssigner(c *gin.Context) {
	c.Set(ginSubContextKey, context.Background())

	c.Next()
}

func defaultNoRouteHandler(c *gin.Context) {
	responder.ErrorWithStatus(http.StatusNotFound, c.Writer, errors.New("the requested resource could not be found"))
}
