package routers

import (
	"encoding/json"
	"fmt"
	"github.com/andresdev99/gambit/db"
	"github.com/andresdev99/gambit/models"
	"github.com/aws/aws-lambda-go/events"
	"strconv"
)

func InsertCategory(body, user string) (int, string) {
	var t models.Category

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error in received data " + err.Error()
	}
	if len(t.CategName) == 0 {
		return 400, "Specify the Category Name"
	}

	if len(t.CategPath) == 0 {
		return 400, "Specify the Category Path"
	}
	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	result, err2 := db.InsertCategory(t)
	if err2 != nil {
		return 400, "Error when inserting Category" + t.CategName + " > " + err2.Error()
	}

	return 200, fmt.Sprintf("{ CategID: %s}", strconv.Itoa(int(result)))
}

func UpdateCategory(body, user string, id int) (int, string) {
	var t models.Category
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return 400, "Error in received data " + err.Error()
	}

	if len(t.CategPath) == 0 && len(t.CategName) == 0 {
		return 400, "Should specify categName and CategName"

	}
	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	t.CategID = id
	err2 := db.UpdateCategory(t)
	if err2 != nil {
		return 400, "Error when trying to Update Category " + strconv.Itoa(id) + " > " + err2.Error()
	}
	return 200, "updated"
}

func DeleteCategory(user string, id int) (int, string) {
	isAdmin, msg := db.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	err2 := db.DeleteCategory(id)
	if err2 != nil {
		return 400, "Error when trying to Delete Category " + strconv.Itoa(id) + " > " + err2.Error()
	}
	return 200, "Deleted"
}

//func GetCategories(body string, request events.APIGatewayV2HTTPRequest) (int, string) {
//	var CategId int
//	var Slug string
//	var err error
//
//	if CategId := request.QueryStringParameters["categId"]; len(CategId) > 0 {
//		CategId, err := strconv.Atoi(CategId)
//
//		if err != nil {
//			return 500, "Error when trying to convert to integer " + strconv.Itoa(CategId)
//		}
//	} else if Slug := request.QueryStringParameters["slug"]; len(Slug) > 0 {
//		// just assign slug
//	}
//
//	list, err2 := db.GetCategories(CategId)
//
//	if err2 != nil {
//		return 400, "Error when retrieving categories " + err2.Error()
//	}
//
//	Categ, err3 := json.Marshal(list)
//
//	if err3 != nil {
//		return 400, "Error when converting categories" + err3.Error()
//	}
//
//	return 200, string(Categ)
//}

func GetCategories(request events.APIGatewayV2HTTPRequest) (int, string) {
	var categID int
	var slug string
	var err error

	queryParams := request.QueryStringParameters

	if rawID := queryParams["categId"]; len(rawID) > 0 {
		categID, err = strconv.Atoi(rawID)
		if err != nil {
			return 400, "Invalid 'categId' parameter: must be an integer"
		}
	} else if s := queryParams["slug"]; len(s) > 0 {
		slug = s
	} else {
		return 400, "Missing required parameter: 'categId' or 'slug'"
	}

	var list []models.Category
	list, err = db.GetCategories(categID, slug)

	if err != nil {
		return 500, "Error retrieving categories: " + err.Error()
	}

	jsonData, err := json.Marshal(list)
	if err != nil {
		return 500, "Error encoding categories to JSON: " + err.Error()
	}

	return 200, string(jsonData)
}
