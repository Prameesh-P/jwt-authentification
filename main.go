package main

import (
	"log"
	"net/http"
)

// starting main func
func main() {
	//three handlers
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/refresh", Refresh)
	//starting port
	log.Fatal(http.ListenAndServe(":8080", nil))
}
