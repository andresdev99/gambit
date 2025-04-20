package awsgo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var Ctx context.Context
var Cfg aws.Config
var err error

func InitializeAWS() {
	Ctx = context.TODO()
	region := "us-east-1"
	Cfg, err = config.LoadDefaultConfig(Ctx, config.WithDefaultRegion(region))

	if err != nil {
		panic("Error when loading configuration .aws/config" + err.Error())
	}
}
