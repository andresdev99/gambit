package routers

import (
	"encoding/json"
	"fmt"
	"github.com/andresdev99/gambit/db"
	"github.com/andresdev99/gambit/models"
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

//func DeleteCategory(user string, id int) (int, string) {
//	isAdmin, msg := db.UserIsAdmin(user)
//	if !isAdmin {
//		return 400, msg
//	}
//
//	err2 := db.DeleteCategory(id)
//	if err2 != nil {
//		return 400, "Error when trying to Delete Category " + strconv.Itoa(id) + " > " + err2.Error()
//	}
//	return 200, "Deleted"
//}
//
//func GetCategories(request events.APIGatewayV2HTTPRequest) (int, string) {
//	var categID int
//	var slug string
//	var err error
//
//	queryParams := request.QueryStringParameters
//
//	if rawID := queryParams["categId"]; len(rawID) > 0 {
//		categID, err = strconv.Atoi(rawID)
//		if err != nil {
//			return 400, "Invalid 'categId' parameter: must be an integer"
//		}
//	} else if s := queryParams["slug"]; len(s) > 0 {
//		slug = s
//	}
//
//	var list []models.Category
//	list, err = db.GetCategories(categID, slug)
//
//	if err != nil {
//		return 500, "Error retrieving categories: " + err.Error()
//	}
//
//	jsonData, err := json.Marshal(list)
//	if err != nil {
//		return 500, "Error encoding categories to JSON: " + err.Error()
//	}
//
//	return 200, string(jsonData)
//}
