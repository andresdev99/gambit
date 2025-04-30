package db

import (
	"database/sql"
	"fmt"
	"github.com/andresdev99/gambit/models"
	"github.com/andresdev99/gambit/tools"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

func InsertProduct(p models.Product) (int64, error) {
	if err := DbConnect(); err != nil {
		return 0, err
	}
	defer Db.Close()

	var (
		fields  []string
		holders []string
		args    []any
	)

	if p.ProdDescription != "" {
		fields = append(fields, "Prod_Description")
		args = append(args, tools.ScapeString(p.ProdDescription))
	}
	if p.ProdPrice > 0 {
		fields = append(fields, "Prod_Price")
		args = append(args, strconv.FormatFloat(p.ProdPrice, 'e', -1, 64))
	}
	if p.ProdPath != "" {
		fields = append(fields, "Prod_Path")
		args = append(args, p.ProdPath)
	}
	if p.ProdCategoryID > 0 {
		fields = append(fields, "Prod_CategoryId")
		args = append(args, strconv.Itoa(p.ProdCategoryID))
	}
	if p.ProdStock > 0 {
		fields = append(fields, "Prod_Stock")
		args = append(args, strconv.Itoa(p.ProdStock))
	}
	fields = append(fields, "Prod_Title")
	args = append(args, tools.ScapeString(p.ProdTitle))

	for range fields {
		holders = append(holders, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO products (%s) VALUES (%s)",
		strings.Join(fields, ", "),
		strings.Join(holders, ", "),
	)
	var debugArgs []string
	for _, a := range args {
		switch v := a.(type) {
		case string:
			// %q will wrap it in quotes and escape if necessary
			debugArgs = append(debugArgs, fmt.Sprintf("%q", v))
		default:
			// all other types just get a plain %v
			debugArgs = append(debugArgs, fmt.Sprintf("%v", v))
		}
	}

	// 2) join them and build the debug SQL
	debug := fmt.Sprintf(
		"INSERT INTO products (%s) VALUES (%s)",
		strings.Join(fields, ", "),
		strings.Join(debugArgs, ", "),
	)

	fmt.Printf("Sentence > %s", debug)

	var result sql.Result

	// 5. exec with args
	result, err := Db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("InsertProduct exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertProduct last id: %w", err)
	}
	return id, nil
}
