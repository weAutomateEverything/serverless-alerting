package main

import (
	"fmt"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	client3 "github.com/weAutomateEverything/serverless-alerting/alert/text/client"
	client2 "github.com/weAutomateEverything/serverless-alerting/telegram/getGroupForChat/client"
	"gopkg.in/telegram-bot-api.v4"
)

type idCommand struct {
}

func newIdCommand() Command {
	return &idCommand{}
}

func (idCommand) CommandIdentifier() string {
	return "id"
}

func (idCommand) CommandDescription() string {
	return "Get the groups ID"
}

func (idCommand) RestrictToAuthorised() bool {
	return false
}

func (idCommand) Show(chat uint32) bool {
	return true
}

func (idCommand) Execute(update tgbotapi.Update) {
	group, err := client2.GetGroupForChat(update.Message.Chat.ID)
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	b := fmt.Sprintf("Your group id is %v",group)
	err = client3.SendMessage(group, b)
	if err != nil {
		client.LogLambdaError(err)
		return
	}

}
