package client

import (
	"fmt"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"net/http"
	"strings"
)

func SendMessage(chat string, message string) error{

	resp, err := http.Post(fmt.Sprintf("%v/alerting/alert/text/%v",common.GetDomain(),chat),"application/text",strings.NewReader(message))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}


func SendError(message string) error {
	return SendMessage(common.GetErrorGroup(),message)
}