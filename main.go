package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var snsClient *sns.Client
var topicArn = "arn:aws:sns:ap-northeast-1:xxxxxxxxxxxx:MyTopic"

// コールドスタート時のみ実行
func init() {
	log.Println("Init function executed (cold start).")

	// AWS SDK の設定ロード
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// SNS クライアント作成
	snsClient = sns.NewFromConfig(cfg)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// リクエストをそのまま SNS に送信する例
	message, _ := json.Marshal(map[string]string{
		"input": req.Body,
	})

	// SNS publish
	_, err := snsClient.Publish(ctx, &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		log.Printf("failed to publish to SNS: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"failed to publish"}`,
		}, nil
	}

	// レスポンス返却
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(`{"message": "Published to SNS"}`),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
