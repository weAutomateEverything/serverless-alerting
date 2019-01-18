package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"github.com/weAutomateEverything/serverless-alerting/telegram/getGroupForChat/api"
	"strconv"
)

var d *dynamodb.DynamoDB

func main() {
	lambda.Start(Handle)
}

func init() {
	c := aws.NewConfig()
	s, err := session.NewSession(c)
	if err != nil {
		client.LogLambdaError(err)
		panic(err)
	}

	d = dynamodb.New(s)
}

func Handle(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {

	chat := request.QueryStringParameters["groupId"]

	var key map[string]*dynamodb.AttributeValue

	for true {
		out, err := d.Scan(&dynamodb.ScanInput{
			ExclusiveStartKey: key,
			TableName:         aws.String("hal"),
			FilterExpression:  aws.String("chat = :c"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":c": {
					S: aws.String(chat),
				},
			},
		})

		if err != nil {
			client.LogLambdaError(err)
			return common.ServerError(err)
		}

		if len(out.Items) > 0 {
			g, err := strconv.ParseInt(*out.Items[0]["groupId"].S, 10, 64)
			if err != nil {
				client.LogLambdaError(err)
				return common.ServerError(err)
			}
			r := api.GetGroupForChatResponse{
				Group: int(g),
			}

			b, err := json.Marshal(&r)
			if err != nil {
				client.LogLambdaError(err)
				return common.ServerError(err)
			}

			response.Body = string(b)
			response.StatusCode = 200
			return response, nil
		}

		key = out.LastEvaluatedKey

	}

	response.StatusCode = 204
	return

}

