package serverless

import (
	"context"
	"sync"

	"github.com/aws/aws-lambda-go/events"

	"github.com/maddiesch/serverless/middleware"
	"github.com/maddiesch/serverless/proxy"
	"github.com/maddiesch/serverless/router"
)

// Application contains the methods needed to create the lambda proxy application
type Application interface {
	Middleware() []middleware.Handler

	Adapter() proxy.Lambda

	Router() router.Router
}

var (
	handlerSetupNonce sync.Once
)

type lambdaHandler struct {
	setup sync.Once
	app   Application
}

func (l *lambdaHandler) fn(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return l.app.Adapter().Proxy(ctx, req)
}

// LambdaHandler is the function passed to `lambda.Start`.
// It runs your code and formats responses.
func LambdaHandler(app Application, setup func(Application)) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handler := &lambdaHandler{}

	return handler.fn
}
