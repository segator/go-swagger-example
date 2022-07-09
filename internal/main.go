package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"go-rest-api/pkg/swagger/server/restapi"
	"go-rest-api/pkg/swagger/server/restapi/operations"
	"io"
	"log"
	"os"
)

var sqsClient *sqs.Client
var queueURL string

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	//Connect to SQS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	queueURL = os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		//log.Fatal("SQS_QUEUE_URL must be set!")
	}
	log.Printf("Setup queue %v to be used.\n", queueURL)
	sqsClient = sqs.NewFromConfig(cfg)

	server.Port = 8080
	api.ApplicationHealthzHandler = operations.ApplicationHealthzHandlerFunc(Health)
	api.PutPublishHandler = operations.PutPublishHandlerFunc(PublishMessage)

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func Health(operations.ApplicationHealthzParams) middleware.Responder {
	return operations.NewApplicationHealthzOK().WithPayload("OK")
}

func PublishMessage(publishParams operations.PutPublishParams) middleware.Responder {
	b, err := io.ReadAll(publishParams.HTTPRequest.Body)
	if err != nil {
		log.Fatalln(err)
	}
	output, err := SendMessage(string(b))
	if err != nil {
		return operations.NewPutPublishServiceUnavailable().WithPayload(err.Error())
	}
	return operations.NewPutPublishOK().WithPayload(*output.MessageId)
}

func SendMessage(msg string) (*sqs.SendMessageOutput, error) {
	output, err := sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    &queueURL,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}
