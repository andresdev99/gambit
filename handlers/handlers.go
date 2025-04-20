package handlers

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

func Handlers(path, method, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Printf("Processing %s > %s", path, method)
	//id := request.PathParameters["id"]
	//idNumber, _ := strconv.Atoi(id)

	return 400, "Invalid Method"
}
