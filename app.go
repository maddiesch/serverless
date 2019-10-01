package serverless

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/maddiesch/serverless/middleware"
)

// LambdaProxy allows lambda requests to be processed by serverless
type LambdaProxy interface {
	Proxy(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// Application contains the methods needed to create the lambda proxy application
type Application interface {
	Adapter() LambdaProxy
}

// GinApp is the default gin application for serverless
type GinApp struct {
	engine  *gin.Engine
	adapter *ginadapter.GinLambda
}

var (
	appSetupNonce     sync.Once
	handlerSetupNonce sync.Once
	appInstance       *GinApp
)

var (
	// App is the current global application
	App Application
)

func init() {
	App = SharedApp()
}

// SharedApp returns the singleton instance of the App.
func SharedApp() *GinApp {
	appSetupNonce.Do(createAppInstance)

	return appInstance
}

// ConfigureGin performs the function passed in.
// The passed function allows access to the underlying gin engine.
func (a *GinApp) ConfigureGin(cfg func(*gin.Engine)) {
	cfg(a.engine)
}

// Adapter returns the lambda proxy
func (a *GinApp) Adapter() LambdaProxy {
	return a.adapter
}

// LambdaHandler is the function passed to `lambda.Start`.
// It runs your code and formats responses.
func LambdaHandler(setup func()) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handlerSetupNonce.Do(setup)

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return App.Adapter().Proxy(req)
	}
}

func createAppInstance() {
	app := gin.New()

	app.Use(gin.Recovery())
	app.Use(middleware.Recovery(recoveryHandler))
	app.Use(middleware.Runtime())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(func(str string) { Log(str) }))

	appInstance = &GinApp{
		engine:  app,
		adapter: ginadapter.New(app),
	}
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
	if IsDebug() {
		meta["ErrorDetails"] = fmt.Sprintf("%v", details)
	}

	return ErrorResponse(c, []*Error{
		&Error{
			Status:      "500",
			Code:        "unknown_error",
			Title:       "Unknown Error",
			Description: "An unknown error occurred. Please try your request again",
			Meta:        &meta,
		},
	})
}
