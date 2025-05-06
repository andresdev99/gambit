package routers

import (
	"encoding/json"
	"fmt"
	"github.com/andresdev99/gambit/db"
	"github.com/andresdev99/gambit/models"
	"github.com/aws/aws-lambda-go/events"
	"strconv"
)

func InsertProduct(body, user string) (int, string) {
	var t models.Product

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error in received data " + err.Error()
	}
	if len(t.ProdTitle) == 0 {
		return 400, "Specify the Prod Title"
	}

	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	result, err2 := db.InsertProduct(t)
	if err2 != nil {
		return 400, "Error when inserting Product" + t.ProdTitle + " > " + err2.Error()
	}

	return 200, fmt.Sprintf("{ ProdId: %s}", strconv.Itoa(int(result)))
}

func UpdateProduct(body, user string, id int) (int, string) {
	var t models.Product
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error in received data " + err.Error()
	}

	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	t.ProdID = id
	err2 := db.UpdateProduct(t)
	if err2 != nil {
		return 400, "Error when trying to Update Product " + strconv.Itoa(id) + " > " + err2.Error()
	}
	return 200, "updated"
}

func DeleteProduct(user string, id int) (int, string) {
	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	err2 := db.DeleteProduct(id)
	if err2 != nil {
		return 400, "Error when trying to Delete Product " + strconv.Itoa(id) + " > " + err2.Error()
	}
	return 200, "Deleted"
}

func GetProducts(request events.APIGatewayV2HTTPRequest) (int, string) {
	var p models.Product
	queryParams := request.QueryStringParameters

	// Parse and validate filters
	if rawID := queryParams["prodId"]; rawID != "" {
		if id, err := strconv.Atoi(rawID); err == nil {
			p.ProdID = id
		} else {
			return 400, "Invalid 'prodId' parameter: must be an integer"
		}
	}

	if rawCateg := queryParams["categId"]; rawCateg != "" {
		p.ProdCategoryID, _ = strconv.Atoi(rawCateg)
	}

	if search := queryParams["search"]; search != "" {
		p.ProdTitle = search // will trigger LIKE
	}

	if slug := queryParams["slug"]; slug != "" {
		p.ProdPath = slug
	}

	// Sorting and pagination
	orderField := queryParams["orderField"]
	orderType := queryParams["orderType"]
	pageSize, _ := strconv.Atoi(queryParams["pageSize"])
	page, _ := strconv.Atoi(queryParams["page"])
	offset := (page - 1) * pageSize

	// Query products
	list, err := db.GetProducts(p, orderField, orderType, pageSize, offset)
	if err != nil {
		return 500, "Error retrieving products: " + err.Error()
	}

	jsonData, err := json.Marshal(list)
	if err != nil {
		return 500, "Error encoding products to JSON: " + err.Error()
	}

	return 200, string(jsonData)
}
