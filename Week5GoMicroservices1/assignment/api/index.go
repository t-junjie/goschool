package main

import (
	"fmt"
	"net/http"
)

//Placeholder page for the API's index page
func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to the REST API")
}
