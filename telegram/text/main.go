package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/weAutomateEverything/halMessageClassification/api"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"strconv"
)

var s *ssm.SSM

func init(){
	c := aws.NewConfig()
	sess, err  := session.NewSession(c)
	if err != nil {
		panic(err)
	}
	s = ssm.New(sess)


}
func main(){
	lambda.Start(Handle)
}

func Handle(request events.APIGatewayProxyRequest)(response events.APIGatewayProxyResponse, err error){
	chat := request.PathParameters["chat"]

	token, err := s.GetParameter(&ssm.GetParameterInput{
		Name:aws.String("telegram-key"),
	})

	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}



	bot,err := tgbotapi.NewBotAPI(*token.Parameter.Value)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	chatId, err := strconv.ParseInt(chat,10,64)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}
	log.Printf("Sending message %v to %v",request.Body,chatId)
	msg := tgbotapi.NewMessage(chatId,request.Body)
	msg.ParseMode = "Markdown"

	_, err = bot.Send(msg)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}
	/*
	t := api.TextEvent{
		MessageID: strconv.Itoa(r.MessageID),
		Message:request.Body,
		Chat: chatId,
	}

	b, err := json.Marshal(t)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}


	resp, err := http.Post("https://api.carddevops.co.za/hal/halMessageClassification","application/json",bytes.NewReader(b))
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	resp.Body.Close()
	*/
	return common.ClientError(200)
}


