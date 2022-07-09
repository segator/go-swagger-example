package e2e

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var appURL = flag.String("appUrl", "http://localhost:8080", "Application URL")
var sqsQueueUrl = flag.String("sqsQueueUrl", "", "SQS Queue where we expect messages been sent")

func TestPublishMessage(t *testing.T) {
	require.NotEmpty(t, *appURL, "test arg appUrl required")
	require.NotEmpty(t, *sqsQueueUrl, "test arg sqsQueueUrl required")
	t.Logf("application Url: %s", *appURL)
	t.Logf("SQS Queue Url: %s", *sqsQueueUrl)

	message := fmt.Sprintf("This is a test %s", guuid.New().String())
	t.Run("Send message and receive from sqs", func(t *testing.T) {
		requestBody := bytes.NewBuffer([]byte(message))
		resp, err := http.Post(fmt.Sprintf("%s/publish", *appURL), "text/plain", requestBody)
		require.NoError(t, err)

		require.Equal(t, 200, resp.StatusCode, "unexpected status code %d", resp.StatusCode)

		responseBodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		sqsMessageIDResponse := string(responseBodyBytes)

		//Connect to SQS
		cfg, err := config.LoadDefaultConfig(context.TODO())
		require.NoError(t, err)

		sqsClient := sqs.NewFromConfig(cfg)
		receiveMessage, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            sqsQueueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     5,
		})
		sqsMessageReceive := receiveMessage.Messages[0]
		assert.Equal(t, sqsMessageIDResponse, *sqsMessageReceive.MessageId)
		assert.Equal(t, message, *sqsMessageReceive.Body)
	})
}
