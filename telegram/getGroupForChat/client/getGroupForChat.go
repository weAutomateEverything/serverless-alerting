package client

import (
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"github.com/weAutomateEverything/serverless-alerting/telegram/getGroupForChat/api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetGroupForChat(chat int64) (group string, err error) {
	requeuest := fmt.Sprintf("%v/alerting/telegram/chat?groupId=%v",common.GetDomain(),chat)
	log.Println(requeuest)
	out, err := http.Get(requeuest)
	var resp api.GetGroupForChatResponse
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	b,err := ioutil.ReadAll(out.Body)
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	log.Println(string(b))
	err = json.Unmarshal(b,&resp)
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	log.Println(resp)
	group = strconv.FormatInt(int64(resp.Group),10)
	return
}
