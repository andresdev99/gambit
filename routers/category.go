package routers

import (
	"encoding/json"
	"fmt"
	"github.com/andresdev99/gambit/db"
	"github.com/andresdev99/gambit/models"
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
