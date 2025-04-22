package db

import (
	"database/sql"
	"fmt"
	"github.com/andresdev99/gambit/models"
	"github.com/andresdev99/gambit/tools"
	_ "github.com/go-sql-driver/mysql"
	"strings"
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

func UpdateCategory(c models.Category) error {
	fmt.Println("Starting Update Category")
	var categoryValues []string
	err := DbConnect()

	if err != nil {
		return err
	}
	defer Db.Close()

	sentence := "UPDATE category SET "

	if catNameVal := tools.ScapeString(c.CategName); len(catNameVal) > 0 {
		categoryValues = append(categoryValues, fmt.Sprintf("Categ_Name = '%s'", catNameVal))
	}

	if catPath := tools.ScapeString(c.CategPath); len(catPath) > 0 {
		categoryValues = append(categoryValues, fmt.Sprintf("Categ_Path= '%s'", catPath))
	}

	if len(categoryValues) == 0 {
		return fmt.Errorf("no values to update")
	}

	sentence += strings.Join(categoryValues, ", ")
	sentence += fmt.Sprintf(" WHERE Categ_Id = %d", c.CategID)

	fmt.Printf("Sentence > %s", sentence)

	var result sql.Result

	result, err = Db.Exec(sentence)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	rowsAffected, err2 := result.RowsAffected()
	if err2 != nil {
		return err2
	}

	if rowsAffected == 0 {
		return fmt.Errorf("error when Updating")
	}

	fmt.Println("Update Category successfully")
	return nil
}

func DeleteCategory(id int) error {
	fmt.Println("Starting Update Category")

	err := DbConnect()

	if err != nil {
		return err
	}
	defer Db.Close()

	sentence := fmt.Sprintf("DELETE FROM category WHERE Categ_Id = %d", id)

	fmt.Printf("Sentence > %s", sentence)

	var result sql.Result

	result, err = Db.Exec(sentence)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	rowsAffected, err2 := result.RowsAffected()
	if err2 != nil {
		return err2
	}

	if rowsAffected == 0 {
		return fmt.Errorf("error when deleting")
	}

	fmt.Println("Delete Category successfully")
	return nil
}

func GetCategories(CategId int, Slug string) ([]models.Category, error) {
	fmt.Println("Starting Select Categories")

	var Categ []models.Category
	err := DbConnect()

	if err != nil {
		return Categ, err
	}
	defer Db.Close()

	sentence := "SELECT Categ_Id, Categ_Name, Categ_Path FROM category"

	if CategId > 0 {
		sentence += fmt.Sprintf(" WHERE Categ_Id = %d", CategId)
	} else if len(Slug) > 0 {
		sentence += fmt.Sprintf(" WHERE Categ_Path LIKE '%%s%'", Slug)
	}
	fmt.Printf("Sentence > %s", sentence)

	var rows *sql.Rows
	rows, err = Db.Query(sentence)

	for rows.Next() {
		var c models.Category
		var categId sql.NullInt32
		var categName sql.NullString
		var categPath sql.NullString

		err := rows.Scan(&categId, &categName, &categPath)

		if err != nil {
			return Categ, err
		}

		c.CategID = int(categId.Int32)
		c.CategName = categName.String
		c.CategPath = categName.String
		Categ = append(Categ, c)
	}

	fmt.Println("Successfully Get Category")

	return Categ, nil
}
