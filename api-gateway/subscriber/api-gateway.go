package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	apigateway "api-gateway/proto/api-gateway"
)

type ApiGateway struct{}

func (e *ApiGateway) Handle(ctx context.Context, msg *apigateway.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *apigateway.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
