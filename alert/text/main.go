package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"log"
	"net/http"
	"strings"
)

var d *dynamodb.DynamoDB

func main(){
	lambda.Start(Handle)
}

func init(){
	config := aws.Config{}
	s,err:= session.NewSession(&config)
	if err != nil {
		client.LogLambdaError(err)
		panic(err)
	}
	d = dynamodb.New(s)
}

func Handle(request events.APIGatewayProxyRequest)(response events.APIGatewayProxyResponse, err error){
	group := request.PathParameters["groupId"]

	log.Printf("%v - %v",group,request.Body)
	i, err := d.GetItem(&dynamodb.GetItemInput{
		TableName:aws.String("hal"),
		Key: map[string]*dynamodb.AttributeValue{
			"groupId":{
				S: aws.String(group),
			},
		},
	})

	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	if i.Item == nil {
		return common.ClientError(412)
	}

	chatStr, ok := i.Item["chat"]
	if !ok {
		return common.ClientError(412)
	}

	resp, err := http.Post(fmt.Sprintf("%v/alerting/telegram/text/%v",common.GetDomain(),*chatStr.S),"application/text",strings.NewReader(request.Body))

	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	return common.ClientError(resp.StatusCode)

}


