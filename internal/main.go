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
	var cfg aws.Config
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID: "aws",
				URL:         awsEndpoint,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}
	queueURL = os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL must be set!")
	}
	log.Printf("Setup queue %v to be used.\n", queueURL)
	sqsClient = sqs.NewFromConfig(cfg)

	server.Port = 8080
	api.ApplicationHealthzHandler = operations.ApplicationHealthzHandlerFunc(Health)
	api.PostPublishHandler = operations.PostPublishHandlerFunc(PublishMessage)

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func Health(operations.ApplicationHealthzParams) middleware.Responder {
	return operations.NewApplicationHealthzOK().WithPayload("OK")
}

func PublishMessage(publishParams operations.PostPublishParams) middleware.Responder {
	b, err := io.ReadAll(publishParams.HTTPRequest.Body)
	if err != nil {
		log.Fatalln(err)
	}
	output, err := SendMessage(string(b))
	if err != nil {
		return operations.NewPostPublishServiceUnavailable().WithPayload(err.Error())
	}
	return operations.NewPostPublishOK().WithPayload(*output.MessageId)
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
