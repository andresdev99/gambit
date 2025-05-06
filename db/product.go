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

	debug := fmt.Sprintf(
		"INSERT INTO products (%s) VALUES (%s)",
		strings.Join(fields, ", "),
		strings.Join(debugArgs, ", "),
	)

	fmt.Printf("Sentence > %s\n", debug)

	var result sql.Result
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

func UpdateStock(p models.Product) error {
	fmt.Println("Starting Update Stock")

	err := DbConnect()

	if err != nil {
		return err
	}
	defer Db.Close()

	sentence := fmt.Sprintf("UPDATE products SET Prod_Stock = Prod_Stock + %v WHERE Prod_id = %d", p.ProdStock, p.ProdID)

	fmt.Printf("Sentence > %s\n", sentence)

	result, err := Db.Exec(sentence)

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

	fmt.Println("Update Stock successfully")
	return nil
}

func UpdateProduct(p models.Product) error {
	fmt.Println("Starting Update Product")
	err := DbConnect()

	if err != nil {
		return err
	}
	defer Db.Close()

	sentence, args, debug, err := tools.BuildInsertUpdateQuery(p, "products", tools.Update)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Sentence > %s\n", debug)

	var result sql.Result

	result, err = Db.Exec(sentence, args...)

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

	fmt.Println("Update Product successfully")
	return nil
}

func DeleteProduct(id int) error {
	fmt.Println("Starting Update Product")

	err := DbConnect()

	if err != nil {
		return err
	}
	defer Db.Close()

	sentence := fmt.Sprintf("DELETE FROM products WHERE Prod_Id = %d", id)

	fmt.Printf("Sentence > %s\n", sentence)

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

	fmt.Println("Delete Product successfully")
	return nil
}

func GetProducts(filter models.Product, orderBy, orderDir string, limit, offset int) ([]models.Product, error) {
	var products []models.Product

	if err := DbConnect(); err != nil {
		return products, err
	}
	defer Db.Close()

	query, args, debug, err := tools.BuildSelectQuery(filter, "products", orderBy, orderDir, limit, offset)
	if err != nil {
		return products, err
	}

	fmt.Println("Query >", debug)

	rows, err := Db.Query(query, args...)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		var prodID, stock, categoryID sql.NullInt32
		var title, desc, createdAt, updatedAt, path sql.NullString
		var price sql.NullFloat64

		err := rows.Scan(&prodID, &title, &desc, &createdAt, &updatedAt, &price, &path, &categoryID, &stock)
		if err != nil {
			return products, err
		}

		p.ProdID = int(prodID.Int32)
		p.ProdTitle = title.String
		p.ProdDescription = desc.String
		p.ProdCreatedAt = createdAt.String
		p.ProdUpdated = updatedAt.String
		p.ProdPrice = price.Float64
		p.ProdPath = path.String
		p.ProdCategoryID = int(categoryID.Int32)
		p.ProdStock = int(stock.Int32)

		products = append(products, p)
	}

	return products, nil
}
