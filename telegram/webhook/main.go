package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/google/uuid"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	client2 "github.com/weAutomateEverything/serverless-alerting/alert/text/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	client3 "github.com/weAutomateEverything/serverless-alerting/telegram/getGroupForChat/client"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
	"strings"
)

var s *ssm.SSM
var d *dynamodb.DynamoDB

type commandCtor func() Command

var commandList = map[string]commandCtor{}

func main() {
	lambda.Start(Handle)
}

func init() {
	c := aws.Config{}
	sess, err := session.NewSession(&c)
	if err != nil {
		client.LogLambdaError(err)
		panic(err)
	}

	s = ssm.New(sess)
	d = dynamodb.New(sess)

	register(newIdCommand)
}

func register(newfunc commandCtor) {
	id := newfunc().CommandIdentifier()
	commandList[strings.ToLower(id)] = newfunc
}

func Handle(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	log.Println(request.Body)
	update := tgbotapi.Update{}
	err = json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		client.LogLambdaError(err)
		return common.ClientError(400)
	}

	token, err := s.GetParameter(&ssm.GetParameterInput{
		Name: aws.String("telegram-key"),
	})

	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	bot, err := tgbotapi.NewBotAPI(*token.Parameter.Value)
	if err != nil {
		client.LogLambdaError(err)
		return common.ServerError(err)
	}

	if update.Message.NewChatMembers != nil {
		for _, user := range *update.Message.NewChatMembers {
			// Looks like the bot has been added to a new group - lets register the details.
			if user.ID == bot.Self.ID {
				chats, err := d.Scan(&dynamodb.ScanInput{
					TableName:        aws.String("hal"),
					FilterExpression: aws.String("chat = :v"),
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":v": {
							S: aws.String(strconv.FormatInt(update.Message.Chat.ID, 10)),
						},
					},
				})

				if err != nil {
					client.LogLambdaError(err)
					return common.ServerError(err)
				}

				id := ""

				if len(chats.Items) == 0 {
					id = strconv.FormatUint(uint64(uuid.New().ID()), 10)
					_, err := d.PutItem(&dynamodb.PutItemInput{
						TableName: aws.String("hal"),
						Item: map[string]*dynamodb.AttributeValue{
							"groupId": {
								S: aws.String(id),
							},
							"chat": {
								S: aws.String(strconv.FormatInt(update.Message.Chat.ID, 10)),
							},
						},
					})
					if err != nil {
						client.LogLambdaError(err)
						return common.ServerError(err)
					}
				} else {
					id = *chats.Items[0]["groupId"].S
				}

				err = client2.SendMessage(id, fmt.Sprintf("The bot has been successfully registered. Your token is %v", id))
				if err != nil {
					client.LogLambdaError(err)
					return common.ServerError(err)
				}

				return common.ClientError(200)
			}
		}
	}

	if strings.HasPrefix(update.Message.Text, "/") {
		tokens := strings.Split(update.Message.Text, "@")
		cmd := tokens[0][1:]
		c, ok := commandList[cmd]
		if !ok {
			group, err := client3.GetGroupForChat(update.Message.Chat.ID)
			if err != nil {
				client.LogLambdaError(err)
				return common.ServerError(err)
			}
			err = client2.SendMessage(group, fmt.Sprintf("No command found for %v", c))
			if err != nil {
				client.LogLambdaError(err)
				return common.ServerError(err)
			}
		} else {
			c().Execute(update)
		}
	}
	return common.ClientError(200)
}
