package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/EmilioCliff/learn-go/go-songapi-with-db/pkg/routes"
	"log"
)

func main(){
	r := mux.NewRouter()
	routes.SetRoutes(r)
	fmt.Println("Starting server at port 8080")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}