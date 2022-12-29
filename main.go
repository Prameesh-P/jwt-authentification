package main

import (
	"github.com/Prameesh-P/jwt-authentification/handlers"
	"log"
	"net/http"
)

// starting main func
func main() {
	//three handlers
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/refresh", handlers.Refresh)
	//starting port
	log.Fatal(http.ListenAndServe(":8080", nil))
}
