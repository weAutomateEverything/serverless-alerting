package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	client2 "github.com/weAutomateEverything/serverless-alerting/alert/text/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
)

func main() {
	lambda.Start(Handle)
}

func Handle(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	m := client.LambdaError{}
	err = json.Unmarshal([]byte(request.Body), &m)
	if err != nil {
		return common.ServerError(err)
	}

	err = client2.SendError(fmt.Sprintf("*%v*\n%v", m.Lambda, m.Error))
	if err != nil {
		return common.ServerError(err)
	}

	return common.ClientError(200)
}
