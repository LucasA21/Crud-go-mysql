package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Server Runing in http://localhost:8080")

	http.ListenAndServe(":8080", nil)

}
