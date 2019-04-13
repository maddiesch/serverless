package serverless

import (
	"context"
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

	app.Use(middleware.Runtime())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(LogMessage))
	app.Use(gin.Recovery())

	appInstance = &App{
		engine:  app,
		adapter: ginadapter.New(app),
	}
}
