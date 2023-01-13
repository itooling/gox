package oth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var auth = "http://localhost:8080/auth"

func Interceptor(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var param map[string]any
		if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
			fmt.Println(err)
		}
		token := r.Header.Get("token")
		if Auth(token) {
			f(w, r)
		} else {
			w.Write([]byte("auth error"))
		}
	}
}

func Auth(token string) bool {
	if token != "" {
		uri := auth + "?token=" + token
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			fmt.Println(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		if string(data) == "success" {
			return true
		}
	}
	return false
}

func HandleTask() {
	zero := http.NewServeMux()
	zero.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		if token == "xxx" {
			w.Write([]byte("success"))
		}
	})
	go http.ListenAndServe(":8080", zero)

	one := http.NewServeMux()
	one.HandleFunc("/one", Interceptor(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("one...")
		w.Header().Set("content-type", "application/json")
		data, err := json.Marshal(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": map[string]interface{}{
				"name": "jack",
				"age":  20,
				"sex":  "male",
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		w.Write(data)
	}))

	go http.ListenAndServe(":8081", one)

	two := http.NewServeMux()
	two.HandleFunc("/two", Interceptor(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("two...")
		w.Header().Set("content-type", "application/json")
		data, err := json.Marshal(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": map[string]interface{}{
				"name": "tom",
				"age":  30,
				"sex":  "female",
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		w.Write(data)
	}))

	go http.ListenAndServe(":8082", two)

	select {}
}
