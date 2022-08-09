package iohandler

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type IO interface {
	Write(path string, object []byte, tags []types.Tag) error
	Read(path string) ([]byte, error)
	List(path string) ([]string, error)
	Count(path string) (int, error)
}

func NewIO(opt string) IO {
	switch opt {
	case "s3":
		return &IOAwsS3{}
	case "local":
		return &IOLocal{}
	case "ssm":
		return &IOAwsSSM{}
	case "secretManager":
		return &IOAwsSecretManager{}
	default:
		fmt.Println("Invalid option")
		return nil
	}
}
