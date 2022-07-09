# go-rest-api

Simple Rest Api application generated using swagger
that can sent messages to an AWS SQS Queue.


## Pre-requisits

- golang
- swagger

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


```
#for development
go run internal/main.go

#for production just run the built binary file in build step
./bin/go-rest-api
```


## Test

### End-To-End
```
go test ./test/e2e --args sqsQueueUrl=$SQS_QUEUE_URL appUrl=$APP_URL
```