package oth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	HttpClient = &http.Client{
		Timeout: time.Second * 60,
	}
)

func OpenFile(path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("file is opening:", f)
}

func UploadFile(url, path string, param map[string]string) error {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	info, err := os.Stat(path)
	if err == os.ErrNotExist {
		return errors.New(fmt.Sprintf("file %s not exist", path))
	}
	if info.IsDir() {
		return errors.New(fmt.Sprintf("path %s is a directory not a file", path))
	}

	fo, err := writer.CreateFormFile("file", info.Name())
	ErrHandle(err)

	for k, v := range param {
		writer.WriteField(k, v)
	}

	fi, err := os.Open(path)
	ErrHandle(err)
	defer fi.Close()

	_, err = io.Copy(fo, fi)
	ErrHandle(err)

	contentType := writer.FormDataContentType()
	writer.Close()

	req, err := http.NewRequest("POST", url, data)
	ErrHandle(err)

	req.Header.Add("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	ErrHandle(err)
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	ErrHandle(err)

	fmt.Println(string(result))

	return nil
}

func FileList(path string) error {
	info, err := ioutil.ReadDir(path)
	for _, f := range info {
		tmp := filepath.Join(path, f.Name())
		if f.IsDir() {
			FileList(tmp)
		} else {
			fmt.Printf("file:%s\n", tmp)
		}
	}
	return err
}

func IsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}
