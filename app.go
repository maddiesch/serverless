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

// App is the serverless application. It handles Lambda requests.
type App struct {
	engine  *gin.Engine
	adapter *ginadapter.GinLambda
}

var (
	appSetupNonce     sync.Once
	handlerSetupNonce sync.Once
	appInstance       *App
)

// SharedApp returns the singleton instance of the App.
func SharedApp() *App {
	appSetupNonce.Do(createAppInstance)

	return appInstance
}

// ConfigureGin performs the function passed in.
// The passed function allows access to the underlying gin engine.
func (a *App) ConfigureGin(cfg func(*gin.Engine)) {
	cfg(a.engine)
}

// LambdaHandler is the function passed to `lambda.Start`.
// It runs your code and formats responses.
func LambdaHandler(setup func()) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handlerSetupNonce.Do(setup)

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return SharedApp().adapter.Proxy(req)
	}
}

func createAppInstance() {
	app := gin.New()

	app.Use(gin.Recovery())
	app.Use(middleware.Recovery(recoveryHandler))
	app.Use(middleware.Runtime())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(LogMessage))

	appInstance = &App{
		engine:  app,
		adapter: ginadapter.New(app),
	}
}

func recoveryHandler(c *gin.Context, e interface{}) error {
	err := e.(error)

	Logf("Recover From Error: %v", err)

	meta := make(map[string]interface{})
	if IsDebug() {
		meta["ErrorDetails"] = fmt.Sprintf("%v", err)
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
