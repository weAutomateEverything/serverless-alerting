package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetDomain() (*string, error){
	return getData("domain")
}

func GetErrorGroup()(*string, error){
	return getData("error_group")
}

func getData(key string) (*string, error){
	c := aws.NewConfig()
	s, err  := session.NewSession(c)
	if err != nil {
		return nil, err
	}

	ss := ssm.New(s)
	out, err := ss.GetParameter(&ssm.GetParameterInput{
		Name:aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Parameter.Value, nil
}
