package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// viewAllCourses sents a HTTP GET request to retrieve all available course information.
func viewAllCourses() {
	url := baseurl
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with the following error: %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
	}
}
