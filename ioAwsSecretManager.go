package iohandler

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type IOAwsSecretManager struct {
	User string `json:"user"`
	Pass string `json:"password"`
}

/* ***************************************************************
* Write a parameter into the parameter store.
*************************************************************** */
func (c *IOAwsSecretManager) Write(path string, object []byte, tags []types.Tag) error {
	return nil
}

/* ***************************************************************
* Read a parameter from the Parameter Store
*************************************************************** */
func (c *IOAwsSecretManager) Read(path string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	secret, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &path,
	})
	if err != nil {
		log.Println(err)
	}
	return []byte(*secret.SecretString), nil
}

/* ***************************************************************
* List all the parameters from the parameter store with the
* past "path"
*************************************************************** */
func (c *IOAwsSecretManager) List(path string) ([]string, error) {
	return nil, nil
}

/* ***************************************************************
* Count the parameters in the Parameter Store in
* the path "path"
*************************************************************** */
func (c *IOAwsSecretManager) Count(path string) (int, error) {
	return 0, nil
}
