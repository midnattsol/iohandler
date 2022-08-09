package iohandler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type IOAwsS3 struct{}

/* ***************************************************************
* Write an object into S3
*************************************************************** */
func (c *IOAwsS3) Write(path string, object []byte, tags []types.Tag) error {
	var t string
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)
	up := manager.NewUploader(client)
	d := strings.Split(path, ":")

	for _, v := range tags {
		t += *v.Key + "=" + *v.Value + ","
	}
	if 1 < len(t) {
		t = t[:len(t)-1]
	}

	out, err := up.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:  &d[0],
		Body:    bytes.NewReader(object),
		Key:     &d[1],
		Tagging: &t,
	})
	if err != nil {
		log.Fatalf("Unable to load file to s3: %s", err)
		return err
	}
	log.Printf("%s has been written.\n", out.Location)
	return nil
}

/* ***************************************************************
* Read an object from S3
*************************************************************** */
func (c *IOAwsS3) Read(path string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	buffer := []byte{}
	client := s3.NewFromConfig(cfg)
	d := strings.Split(path, ":")

	obj, err := client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &d[0],
			Key:    &d[1],
		})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	obj.Body.Read(buffer)
	return buffer, nil
}

/* ***************************************************************
* List the objects with the prefix "path" from S3
*************************************************************** */
func (c *IOAwsS3) List(path string) ([]string, error) {
	list_files := []string{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	d := strings.Split(path, ":")

	list, err := client.ListObjectsV2(
		context.Background(),
		&s3.ListObjectsV2Input{
			Bucket:  &d[0],
			Prefix:  &d[1],
			MaxKeys: 1000,
		})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, o := range list.Contents {
		list_files = append(list_files, *o.Key)
	}
	return list_files, nil
}

/* ***************************************************************
* Count the objects with the prefix "path" from S3
*************************************************************** */
func (c *IOAwsS3) Count(path string) (int, error) {
	list, err := c.List(path)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return len(list), nil
}
