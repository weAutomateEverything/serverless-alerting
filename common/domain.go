package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetDomain() string {
	return getData("domain")
}

func GetErrorGroup() string {
	return getData("error_group")
}

func getData(key string) string {
	c := aws.NewConfig()
	s, err  := session.NewSession(c)
	if err != nil {
		panic(err)
	}

	ss := ssm.New(s)
	out, err := ss.GetParameter(&ssm.GetParameterInput{
		Name:aws.String(key),
	})
	if err != nil {
		panic(err)
	}
	return *out.Parameter.Value
}
