package ginapp

import (
	"context"
	"net/http"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/maddiesch/serverless/proxy"
)

type ginProxy struct {
	core.RequestAccessor

	proxy.Lambda

	engine *gin.Engine
}

func (p *ginProxy) Proxy(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ginRequest, err := p.ProxyEventToHTTPRequest(req)

	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	p.engine.ServeHTTP(http.ResponseWriter(respWriter), ginRequest)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
