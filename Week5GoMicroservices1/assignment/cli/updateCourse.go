package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

// updateCourse allows users to update course information for a course.
func updateCourse() {
	fmt.Printf(courseCodePrompt)
	// courseCode is assumed to be a string without spaces in between
	var courseCode string
	fmt.Scanln(&courseCode)

	if !valid.IsAlphanumeric(courseCode) {
		fmt.Println(invalidInputStr)
		return
	}

	fmt.Printf(newCourseDetailsPrompt)
	fmt.Println(separator)
	fmt.Printf(courseCodePrompt)
	// courseCode is assumed to be a string without spaces in between
	var newCourseCode string
	fmt.Scanln(&newCourseCode)

	if !valid.IsAlphanumeric(newCourseCode) {
		fmt.Println(invalidInputStr)
		return
	}

	fmt.Printf(courseNamePrompt)
	in := bufio.NewReader(os.Stdin)
	newCourseName, err := in.ReadString('\n')
	//trim newline character
	newCourseName = strings.TrimSuffix(newCourseName, "\n")
	if err != nil {
		fmt.Printf("Course Name could not be read: %s\n", err)
	}

	// remove spaces to validate input
	stripStr := strings.ReplaceAll(newCourseName, " ", "")
	if !valid.IsAlphanumeric(stripStr) {
		fmt.Println(invalidInputStr)
		return
	}

	jsonData := map[string]string{"newCourseName": newCourseName, "newCourseCode": newCourseCode}
	jsonValue, _ := json.Marshal(jsonData)

	url := baseurl + "/" + courseCode
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	req.Header.Set("Content-Type", contentType)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(string(data))
}
