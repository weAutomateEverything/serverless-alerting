package client

import (
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/serverless-alerting/alert/lambda/client"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"github.com/weAutomateEverything/serverless-alerting/telegram/getGroupForChat/api"
	"net/http"
	"strconv"
)

func GetGroupForChat(chat int64) (group string, err error) {
	requeuest := fmt.Sprintf("%v/telegram/chat?groupId=%v",common.GetDomain(),chat)
	out, err := http.Get(requeuest)
	var resp api.GetGroupForChatResponse
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	err = json.NewDecoder(out.Body).Decode(&resp)
	if err != nil {
		client.LogLambdaError(err)
		return
	}
	group = strconv.FormatInt(int64(resp.Group),64)
	return
}
