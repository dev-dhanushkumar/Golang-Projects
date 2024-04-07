package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func setStaticFolder(route *mux.Router) {
	fs := http.FileServer(http.Dir("./public/"))
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
}

func addApproutes(route *mux.Router) {
	setStaticFolder((route))

	route.HandleFunc("/", renderHome)

	route.HandleFunc("/user", getUser).Methods("GET")

	route.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")

	route.HandleFunc("/user", updateUser).Methods("PUT")

	route.HandleFunc("/user", insertUser).Methods("POST")

	fmt.Println("Routes Setup Successfully!")
}
