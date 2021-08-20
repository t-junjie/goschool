package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	valid "github.com/asaskevich/govalidator"
)

// viewCourse sents a HTTP GET request to retrieve all available course information.
func viewCourse() {
	fmt.Printf(courseCodePrompt)
	var courseCode string
	fmt.Scanln(&courseCode)

	if !valid.IsAlphanumeric(courseCode) {
		fmt.Println(invalidInputStr)
		return
	}

	fmt.Println(separator)

	url := baseurl + "/" + courseCode
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("The HTTP request failed with the following error: %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
	}
}
