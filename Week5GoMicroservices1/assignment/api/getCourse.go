package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// getCourse allows users to get information on a particular course
func getCourse(db *sql.DB, res http.ResponseWriter, courseCode string) {
	query := "SELECT * FROM goms_db.Courses WHERE CourseCode = ?"
	result, err := db.Query(query, courseCode)
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
