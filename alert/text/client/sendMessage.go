package client

import (
	"fmt"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"net/http"
	"strings"
)

func SendMessage(chat string, message string) error{
	d, err := common.GetDomain()
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%v/telegram/alert/%v",*d,chat),"application/text",strings.NewReader(message))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}


func SendError(message string) error {
	e, err := common.GetErrorGroup()
	if err != nil {
		return err
	}

	return SendMessage(*e,message)
}