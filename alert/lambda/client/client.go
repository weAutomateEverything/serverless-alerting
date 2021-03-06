package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/telegram"
	"github.com/weAutomateEverything/serverless-alerting/common"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func LogLambdaError(err error) {
	msg := fmt.Sprintf("ERROR: %v\n %v", telegram.Escape(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")), telegram.Escape(err.Error()))
	log.Printf(msg)
	l := LambdaError{
		Error:  err.Error(),
		Lambda: os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
	}

	b, err := json.Marshal(&l)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(fmt.Sprintf("https://%v/alerting/alert/lambda", common.GetDomain()), "application/text", bytes.NewReader(b))
	if err == nil {
		resp.Body.Close()
	} else {
		log.Printf("HAL ERROR: %v", err.Error())
	}
	debug.PrintStack()
}

type LambdaError struct {
	Error  string
	Lambda string
}
