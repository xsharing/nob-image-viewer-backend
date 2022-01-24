package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	imageviewerapi "github.com/image-viewer/nob-image-viewer-backend"
	imageviewer "github.com/image-viewer/nob-image-viewer-backend/gen/image_viewer"
	goahttp "goa.design/goa/v3/http"
)

var logger *log.Logger
var adapter *GoaLambdaV2

type GoaLambdaV2 struct {
	core.RequestAccessorV2

	goaEngine goahttp.Muxer
}

func NewV2(goa *goahttp.Muxer) *GoaLambdaV2 {
	return &GoaLambdaV2{goaEngine: *goa}
}

func init() {
	// Setup logger. Replace logger with your own log package of choice.
	logger = log.New(os.Stderr, "[imageviewerapi] ", log.Ltime)

	// Initialize the services.
	var (
		imageViewerSvc imageviewer.Service
	)
	{
		imageViewerSvc = imageviewerapi.NewImageViewer(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		imageViewerEndpoints *imageviewer.Endpoints
	)
	{
		imageViewerEndpoints = imageviewer.NewEndpoints(imageViewerSvc)
	}

	mux := createHTTPServer(imageViewerEndpoints, logger, false)

	adapter = NewV2(mux)

	logger.Print("mux created!")
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	request, err := adapter.EventToRequestWithContext(ctx, req)
	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}
	respWriter := core.NewProxyResponseWriterV2()

	adapter.goaEngine.ServeHTTP(http.ResponseWriter(respWriter), request)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeoutV2(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}

func main() {
	logger.Printf("lambda will start")
	lambda.Start(Handler)
}
