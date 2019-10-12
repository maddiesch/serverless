package proxy

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// Lambda allows AWS Lambda requests to be processed by serverless
type Lambda interface {
	Proxy(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
