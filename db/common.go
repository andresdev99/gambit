package db

import (
	"database/sql"
	"fmt"
	"github.com/andresdev99/gambit/models"
	"github.com/andresdev99/gambit/secretm"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var SecretModel models.SecretRDSJson
var err error
var Db *sql.DB

func ReadSecret() error {
	SecretModel, err = secretm.GetSecret(os.Getenv("SecretName"))
	return err
}

func DbConnect() error {
	Db, err = sql.Open("mysql", ConnStr(SecretModel))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = Db.Ping()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Connected Successfully")
	return nil
}

func ConnStr(keys models.SecretRDSJson) string {
	var dbUser, authToken, dbEndpoint, dbName string

	dbUser = keys.UserName
	authToken = keys.Password
	dbEndpoint = keys.Host
	dbName = "gambit"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true", dbUser, authToken, dbEndpoint, dbName)
	fmt.Println(dsn)
	return dsn
}

func UserIsAdmin(userUUID string) (bool, string) {
	fmt.Println("Starts User Is Admin Func")

	err := DbConnect()

	if err != nil {
		return false, err.Error()
	}

	defer Db.Close()

	sentence := fmt.Sprintf("SELECT 1 FROM users where User_UUID ='%s' AND User_Status = 0", userUUID)
	fmt.Printf("Sentence > %s\n", sentence)

	rows, err := Db.Query(sentence)
	if err != nil {
		return false, err.Error()
	}

	var value string

	rows.Next()
	err = rows.Scan(&value)
	//if err != nil {
	//	return false, err.Error()
	//}

	fmt.Println("User Is Admin > Execution successfully")

	if value == "1" {
		return true, ""
	}

	return false, "User is not Admin > " + userUUID
}
