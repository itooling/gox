package oth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func FileServer(path string) error {
	fmt.Println("starting server...")
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(path))))
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir(path))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}

func FileUploadServer() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		fmt.Println(token)

		file, header, err := r.FormFile("file")
		ErrHandle(err)
		defer file.Close()

		filename := header.Filename

		fi := bufio.NewReader(file)

		fo, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
		ErrHandle(err)
		defer fo.Close()

		_, err = io.Copy(fo, fi)
		ErrHandle(err)

	})

	http.ListenAndServe(":8080", nil)
}

/*
生成 ca.key 私钥
openssl genrsa -out ca.key 4096

制作解密后的私钥（一般无此必要）
openssl rsa -in ca.key -out ca_decrypted.key

生成 ca.crt 根证书（公钥）：
openssl req -new -x509 -days 7304 -key ca.key -out ca.crt

制作生成网站的证书并用签名认证
如果证书已经存在或者想沿用，则可以直接从这一步开始。

生成证书私钥：
openssl genrsa -out server.pem 4096

制作解密后的证书私钥：
openssl rsa -in server.pem -out server.key

生成签名请求：
openssl req -new -key server.pem -out server.csr

在 common name 中填入网站域名，如 www.baidu.com 即可生成该站点的证书，同时也可以使用泛域名如 *.baidu.com 来生成所有二级域名可用的网站证书。

用证书进行签名：
openssl ca -policy policy_anything -days 3652 -cert ca.crt -keyfile ca.key -in server.csr -out server.crt
*/
func HttpServer(scheme string) {
	var c chan int
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var tls string
		if r.TLS == nil {
			tls = "http"
		} else {
			tls = "https"
		}
		fmt.Printf("scheme is: %s\n", tls)
		w.Header().Set("content-type", "application/json")
		data, _ := json.Marshal(map[string]interface{}{"code": 200, "msg": "hello world"})
		w.Write(data)
	})
	if scheme == "http" {
		go func() {
			err := http.ListenAndServe(":80", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	} else if scheme == "https" {
		go func() {
			err := http.ListenAndServeTLS(":443", "pem/ca.crt", "pem/ca.key", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	} else if scheme == "all" {
		go func() {
			err := http.ListenAndServe(":80", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
		go func() {
			err := http.ListenAndServeTLS(":443", "pem/ca.crt", "pem/ca.key", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	fmt.Println("server running...")
	select {
	case <-c:
		fmt.Println("server shutdown...")
	}
}

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello world"))
		rw.WriteHeader(200)
	})
	http.ListenAndServe(":8080", mux)
}
