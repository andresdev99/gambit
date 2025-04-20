package main

import (
	"context"
	"github.com/andresdev99/gambit/awsgo"
	"github.com/andresdev99/gambit/db"
	"github.com/andresdev99/gambit/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"
)

func main() {
	lambda.Start(ExecuteLambda)
}

func ExecuteLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	awsgo.InitializeAWS()
	if !ValidateParameters() {
		panic("Parameters Error: should send SecretName, UserPoolId, Region and UrlPrefix")
	}
	var res *events.APIGatewayProxyResponse
	path := strings.Replace(request.RawPath, os.Getenv("UrlPrefix"), "", -1)
	method := request.RequestContext.HTTP.Method
	body := request.Body
	headers := request.Headers

	if err := db.ReadSecret(); err != nil {
		panic("Error when reading Secret")
	}

	status, message := handlers.Handlers(path, method, body, headers, request)

	headersResponse := map[string]string{
		"content-Type": "application/json",
	}

	res = &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
		Headers:    headersResponse,
	}
	return res, nil
}

func ValidateParameters() bool {
	requiredVars := []string{"SecretName", "UrlPrefix"}

	for _, envParam := range requiredVars {
		if _, ok := os.LookupEnv(envParam); !ok {
			return false
		}
	}
	return true
}
