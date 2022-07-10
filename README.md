# go-rest-api

Simple Rest Api application generated using swagger
that can sent messages to an AWS SQS Queue.


## Pre-requisits

- https://go.dev/ >= 1.17
- https://github.com/go-swagger/go-swagger

## API definition
Api endpoints are generated automatically using swagger.

[Application API definition](./pkg/swagger/swagger.yml)

```
#Validate swagger definition
swagger validate pkg/swagger/swagger.yml
```
## Build
```
#Generate code from swagger definition 
go generate go-rest-api/internal go-rest-api/pkg/swagger
go build -o bin/go-rest-api internal/main.go
```

## Run the app

### Configure
The application needs an AWS Account configured to be able to send messages to an SQS Queue.

All the possible authentication methods are described in https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

The application automatically supports:
- `$HOME/.aws/credentials`
- Environment variables

AWS Endpoint can be override with
env var `AWS_ENDPOINT=http://xxxxxx`


### Run
```
#for development
go run internal/main.go

#for production just run the built binary file in build step
./bin/go-rest-api

```


## Test

### End-To-End

#### Prerequisites
- AWS SQS Queue
- This App running

#### Run the test
```
go test ./test/e2e --sqsQueueUrl=$SQS_QUEUE_URL --appUrl=$APP_URL
```