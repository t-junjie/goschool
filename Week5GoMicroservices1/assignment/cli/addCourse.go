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

// addCourse allows users to add a course via the REST API.
func addCourse() {
	fmt.Printf(courseCodePrompt)
	// courseCode is assumed to be a string without spaces in between
	var courseCode string
	fmt.Scanln(&courseCode)
	if !valid.IsAlphanumeric(courseCode) {
		fmt.Println(invalidInputStr)
		return
	}

	fmt.Printf(courseNamePrompt)
	in := bufio.NewReader(os.Stdin)
	courseName, err := in.ReadString('\n')
	//trim newline character
	courseName = strings.TrimSuffix(courseName, "\n")
	if err != nil {
		fmt.Printf("Course Name could not be read: %s\n", err)
	}

	// remove spaces to validate input
	stripStr := strings.ReplaceAll(courseName, " ", "")
	if !valid.IsAlphanumeric(stripStr) {
		fmt.Println(invalidInputStr)
		return
	}

	jsonData := map[string]string{"courseName": courseName}
	jsonValue, _ := json.Marshal(jsonData)

	res, err := http.Post(baseurl+"/"+courseCode, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(string(data))
}
