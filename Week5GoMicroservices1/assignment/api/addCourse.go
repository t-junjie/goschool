package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// addCourse allows users to create a new course
func addCourse(db *sql.DB, res http.ResponseWriter, req *http.Request, courseCode string) {
	var newCourse courseInfo
	var responseString string
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseString = "HTTP 422: Please supply course information in JSON format."
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(responseString))
	}

	json.Unmarshal(reqBody, &newCourse)

	if newCourse.CourseName == "" {
		responseString = "HTTP 422: Course name cannot be an empty string."
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(responseString))
		return
	}

	// check if course already exists
	existQuery := "SELECT CourseCode FROM Courses WHERE CourseCode =?"
	exists := rowExists(existQuery, courseCode)

	// insert only if course does not exist
	if !exists {
		query := "INSERT INTO goms_db.Courses (CourseCode, CourseName) VALUES(?,?)"
		results, err := db.Exec(query, courseCode, newCourse.CourseName)
		if err != nil {
			panic(err)
		} else {
			rows, _ := results.RowsAffected()
			fmt.Printf("Rows affected: %v\n", rows)
			fmt.Printf("Course: %s - %s has been added to the database.\n", courseCode, newCourse.CourseName)
		}

		responseString = fmt.Sprintf("%s: %s was successfully added.\n", courseCode, newCourse.CourseName)
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(responseString))
	} else {
		responseString = "HTTP 409: Duplicate course code."
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(responseString))
	}

}
