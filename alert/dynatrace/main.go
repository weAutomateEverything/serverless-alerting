package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	client2 "github.com/weAutomateEverything/serverless-alerting/alert/text/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
)

func main(){
	lambda.Start(Handle)
}

func Handle(request events.APIGatewayProxyRequest)(response events.APIGatewayProxyResponse, err error ){
	var m msg
	err = json.Unmarshal([]byte(request.Body),&m)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	err = client2.SendError(m.Alert)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	return common.ClientError(200)
}

type msg struct {
	Alert string `json:"alert"`
}


