package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Server will started at http://localhost:8080/")

	connecDatabase()

	route := mux.NewRouter()

	addApproutes(route)

	log.Fatal(http.ListenAndServe(":8080", route))
}
