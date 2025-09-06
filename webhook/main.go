package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// SNSから送られてくるメッセージ構造体
type SNSMessage struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`

	// SubscriptionConfirmation 時のみ
	Token        string `json:"Token"`
	SubscribeURL string `json:"SubscribeURL"`
}

// Lambda Handler
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var msg SNSMessage
	if err := json.Unmarshal([]byte(req.Body), &msg); err != nil {
		fmt.Printf("failed to parse SNS message: %v\n", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	switch msg.Type {
	case "SubscriptionConfirmation":
		// 購読確認 → SubscribeURL を叩いて確定
		fmt.Printf("Subscription confirmation received. URL=%s\n", msg.SubscribeURL)
		resp, err := http.Get(msg.SubscribeURL)
		if err != nil {
			fmt.Printf("failed to confirm subscription: %v\n", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Subscription confirmed response: %s\n", string(body))

	case "Notification":
		// 通知受信 → メッセージを処理
		fmt.Printf("Notification received. Message=%s\n", msg.Message)

	case "UnsubscribeConfirmation":
		fmt.Printf("Unsubscribe confirmation: %s\n", msg.Message)

	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "OK"}, nil
}

func main() {
	lambda.Start(handler)
}
