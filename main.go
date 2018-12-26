package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/answers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			{
				var httpHeaders = w.Header()
				httpHeaders.Add("Access-Control-Allow-Origin", "*")
				fmt.Fprint(w, "")
			}
		case "GET":
			{
				body, err := ioutil.ReadFile("answers.json")
				if err != nil {
					panic(err)
				}
				var httpHeaders = w.Header()
				httpHeaders.Add("Access-Control-Allow-Origin", "*")
				log.Println(string(body))
				fmt.Fprint(w, string(body))
			}
		}
	})

	router.HandleFunc("/answers/save", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			{
				var httpHeaders = w.Header()
				httpHeaders.Add("Access-Control-Allow-Origin", "*")
				fmt.Fprint(w, "")
			}
		case "POST":
			{
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				err1 := ioutil.WriteFile("answers.json", data, 0644)
				if err1 != nil {
					panic(err1)
				}
				var httpHeaders = w.Header()
				httpHeaders.Add("Access-Control-Allow-Origin", "*")
			}
		}
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var httpHeaders = w.Header()
		httpHeaders.Add("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "")
	})

	http.Handle("/", router)

	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}

}
