package secretm

import (
	"encoding/json"
	"fmt"
	"github.com/andresdev99/gambit/awsgo"
	"github.com/andresdev99/gambit/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(secretName string) (models.SecretRDSJson, error) {
	var secretData models.SecretRDSJson
	fmt.Println(" > Get Secret " + secretName)
	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		fmt.Println(err)
		return secretData, err
	}
	err = json.Unmarshal([]byte(*key.SecretString), &secretData)
	if err != nil {
		return secretData, err
	}
	fmt.Println(" > Reading Secret Ok " + secretName)
	return secretData, nil
}
