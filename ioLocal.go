package iohandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type IOLocal struct{}

/* ***************************************************************
* Write a file in the path "dest".
* It creates the folders in case they don't exist
*************************************************************** */
func (c *IOLocal) Write(path string, object []byte, tag []types.Tag) error {
	splitted := strings.Split(path, "/")[:]
	p := strings.Join(splitted[:len(splitted)-1], "/")
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = ioutil.WriteFile(path, object, 0650)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/* ***************************************************************
* Read a file from the destiny "path"
*************************************************************** */
func (c *IOLocal) Read(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

/* ***************************************************************
* List the elements in the directory "path"
*************************************************************** */
func (c *IOLocal) List(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	list_files := []string{}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, f := range files {
		list_files = append(list_files, f.Name())
	}
	return list_files, nil
}

/* ***************************************************************
* Count the number of element in the path "path"
*************************************************************** */
func (c *IOLocal) Count(path string) (int, error) {
	list, err := c.List(path)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return len(list), nil
}
