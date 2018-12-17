package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/answers", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var httpHeaders = w.Header()
		httpHeaders.Add("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
