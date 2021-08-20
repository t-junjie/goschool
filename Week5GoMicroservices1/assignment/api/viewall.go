package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// viewall retrieves all course information from database
func viewall(res http.ResponseWriter, req *http.Request) {
	result, err := db.Query("SELECT * FROM goms_db.Courses")
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var course Course
		err = result.Scan(&course.ID, &course.CourseCode, &course.CourseName)
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Course: %s - %s has been retrieved from the database.\n", course.CourseCode, course.CourseName)
		courseInfo := fmt.Sprintf("Course ID: %d, Course Code: %s, Course Title: %s", course.ID, course.CourseCode, course.CourseName)
		json.NewEncoder(res).Encode(courseInfo)
	}
}
