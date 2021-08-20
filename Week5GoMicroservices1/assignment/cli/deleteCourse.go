package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	valid "github.com/asaskevich/govalidator"
)

//delCourse allows users to delete a course
func deleteCourse() {
	fmt.Printf(courseCodePrompt)
	// courseCode is assumed to be a string without spaces in between
	var courseCode string
	fmt.Scanln(&courseCode)
	if !valid.IsAlphanumeric(courseCode) {
		fmt.Println(invalidInputStr)
		return
	}

	url := baseurl + "/" + courseCode
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(string(data))
}
