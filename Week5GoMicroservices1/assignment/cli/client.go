// Package client is a REST API client that allows users to perform CRUD operations on courses
package main

import (
	"fmt"
	"time"

	valid "github.com/asaskevich/govalidator"
)

const (
	baseurl                = "http://localhost:5221/api/v1/courses"
	separator              = "====================================="
	contentType            = "application/json"
	choosePrompt           = "Select an option: "
	newCourseDetailsPrompt = "Please key in new information for the course.\nPress enter to retain the value for the field.\n"
	courseCodePrompt       = "Enter course code: "
	courseNamePrompt       = "Enter course name: "
	invalidInputStr        = "Invalid Input. Input must be alphanumeric."
	mainPrompt             = `
=====================================
Welcome to the Course API Client!

1. View All Courses
2. Add Course
3. View Course
4. Update Course
5. Delete Course
=====================================
`
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

func main() {
	for {
		chooseOptions()
	}
}

// chooseOptions is a menu to allow users to choose which CRUD operation to perform on the courses.
func chooseOptions() {
	fmt.Println(mainPrompt)
	fmt.Println(choosePrompt)

	var input int
	fmt.Scanln(&input)

	switch input {
	case 1:
		fmt.Println(separator)
		viewAllCourses()
		fmt.Println(separator)
	case 2:
		fmt.Println(separator)
		addCourse()
		fmt.Println(separator)
	case 3:
		fmt.Println(separator)
		viewCourse()
		fmt.Println(separator)
	case 4:
		fmt.Println(separator)
		updateCourse()
		fmt.Println(separator)
	case 5:
		fmt.Println(separator)
		deleteCourse()
		fmt.Println(separator)
	default:
		fmt.Println(separator)
		fmt.Println("Invalid Input. Please try again.")
		fmt.Println(separator)
	}

	time.Sleep(1500 * time.Millisecond)
}
