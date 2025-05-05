package handlers

import (
	"fmt"
	"github.com/andresdev99/gambit/auth"
	"github.com/andresdev99/gambit/routers"
	"github.com/aws/aws-lambda-go/events"
	"strconv"
)

func Handlers(path, method, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Printf("Processing %s > %s", path, method)
	id := request.PathParameters["id"]
	idNumber, _ := strconv.Atoi(id)
	isOk, statusCode, user := validateAuthorization(path, method, headers)

	if !isOk {
		return statusCode, user
	}

	switch path[1:5] {
	case "user":
		return UsersProcess(body, path, method, user, id, request)
	case "prod":
		return ProductsProcess(body, path, method, user, idNumber, request)
	case "stoc":
		return StockProcess(body, path, method, user, idNumber, request)
	case "addr":
		return AddressProcess(body, path, method, user, idNumber, request)
	case "cate":
		return CategoryProcess(body, path, method, user, idNumber, request)
	case "orde":
		return OrderProcess(body, path, method, user, idNumber, request)
	}

	return 400, "Invalid Method for path > " + path + " " + path[1:5]
}

func validateAuthorization(path, method string, headers map[string]string) (bool, int, string) {
	if (path == "product" || path == "category") && method == "GET" {
		return true, 200, ""
	}

	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token is required"
	}

	ok, err, msg := auth.ValidateToken(token)
	if !ok {
		tokenErrMsg := "Token error: "
		if err != nil {
			fmt.Println(tokenErrMsg, err.Error())
			return false, 401, err.Error()
		}
		fmt.Println(tokenErrMsg, msg)
		return false, 401, msg
	}
	fmt.Println("Token OK")
	return true, 200, msg
}

func UsersProcess(body, path, method, user, id string, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Invalid Method For UsersProcess"
}

func ProductsProcess(body, path, method, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	switch method {
	case "POST":
		return routers.InsertProduct(body, user)
	case "PUT":
		return routers.UpdateProduct(body, user, id)
	}
	return 400, "Invalid Method For ProductsProcess"
}

func CategoryProcess(body, path, method, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	switch method {
	case "POST":
		return routers.InsertCategory(body, user)
	case "PUT":
		return routers.UpdateCategory(body, user, id)
	case "DELETE":
		return routers.DeleteCategory(user, id)
	case "GET":
		return routers.GetCategories(request)
	}
	return 400, "Invalid Method For Category"
}

func StockProcess(body, path, method, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Invalid Method For StockProcess"
}

func AddressProcess(body, path, method, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Invalid Method For AddressProcess"
}

func OrderProcess(body, path, method, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {

	return 400, "Invalid Method For OrderProcess"
}
