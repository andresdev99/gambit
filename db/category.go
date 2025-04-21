package db

import (
	"database/sql"
	"fmt"
	"github.com/andresdev99/gambit/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertCategory(c models.Category) (int64, error) {
	fmt.Println("Starting Insert Category")

	err := DbConnect()

	if err != nil {
		return 0, err
	}
	defer Db.Close()
	sentence := fmt.Sprintf("INSERT INTO category (Categ_Name, Categ_Path) VALUES ('%s', '%s')", c.CategName, c.CategPath)
	fmt.Printf("Sentence > %s", sentence)

	var result sql.Result

	result, err = Db.Exec(sentence)

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	lastInsertedId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}
	fmt.Println("Insert Category successfully")
	return lastInsertedId, nil
}
