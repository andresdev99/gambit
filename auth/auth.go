package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenJSON struct {
	Sub       string
	Event_Id  string
	Token_use string
	Scope     string
	Auth_time int
	Iss       string
	Exp       int
	Iat       int
	Client_id string
	Username  string
}

func ValidateToken(token string) (bool, error, string) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		message := "Not valid Token"
		fmt.Println(message)
		return false, nil, message
	}

	userInfo, err := base64.StdEncoding.DecodeString(parts[1])

	if err != nil {
		fmt.Println("Cannot be Decoded the Token Part: ", err.Error())
		return false, err, err.Error()
	}

	var tkj TokenJSON
	err = json.Unmarshal(userInfo, &tkj)
	if err != nil {
		fmt.Println("Cannot be Decoded the JSON structure ", err.Error())
		return false, err, err.Error()
	}

	now := time.Now()
	tm := time.Unix(int64(tkj.Exp), 0)

	if tm.Before(now) {
		fmt.Printf("Expiration date = %s\n", tm.String())
		message := "Expired Token"
		fmt.Println(message)
		return false, nil, message
	}

	return true, nil, string(tkj.Username)
}
