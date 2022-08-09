package iohandler

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type IOAwsSSM struct{}

/* ***************************************************************
* Write a parameter into the parameter store.
*************************************************************** */
func (c *IOAwsSSM) Write(path string, object []byte, tags []types.Tag) error {
	var ss string = string(object)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return err
	}

	client := ssm.NewFromConfig(cfg)
	client.PutParameter(context.Background(), &ssm.PutParameterInput{
		Name:  &path,
		Value: &ss,
		Tags:  tags,
	})
	return nil
}

/* ***************************************************************
* Read a parameter from the Parameter Store
*************************************************************** */
func (c *IOAwsSSM) Read(path string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)
	param, err := client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           &path,
		WithDecryption: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	return []byte(*param.Parameter.Value), nil
}

/* ***************************************************************
* List all the parameters from the parameter store with the
* past "path"
*************************************************************** */
func (c *IOAwsSSM) List(path string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	var output []string
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)
	params, err := client.DescribeParameters(context.Background(), &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{
			{
				Key:    aws.String("Path"),
				Values: []string{path},
			},
		},
	})
	for _, v := range params.Parameters {
		output = append(output, *v.Name)
	}
	return output, err
}

/* ***************************************************************
* Count the parameters in the Parameter Store in
* the path "path"
*************************************************************** */
func (c *IOAwsSSM) Count(path string) (int, error) {
	p, err := c.List(path)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return len(p), err
}
