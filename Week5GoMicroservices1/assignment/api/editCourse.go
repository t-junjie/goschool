package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// editCourse allows users to update course information for a course.
func editCourse(db *sql.DB, res http.ResponseWriter, req *http.Request, courseCode string) {
	var newCourse newCourseInfo
	var responseString string
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseString = "HTTP 422: Please supply course information in JSON format."
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(responseString))
	}

	// check if course already exists
	existQuery := "SELECT CourseCode FROM Courses WHERE CourseCode =?"
	exists := rowExists(existQuery, courseCode)

	if !exists {
		responseString = "HTTP 404: " + courseCode + " does not exist. Cannot edit a non-existant course.\n"
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(responseString))
		return
	}

	var currentCourse Course
	query := "SELECT * FROM goms_db.Courses WHERE CourseCode = ?"
	result, err := db.Query(query, courseCode)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&currentCourse.ID, &currentCourse.CourseCode, &currentCourse.CourseName)
		if err != nil {
			panic(err.Error())
		}
	}

	json.Unmarshal(reqBody, &newCourse)
	if newCourse.NewCourseCode == "" {
		newCourse.NewCourseCode = currentCourse.CourseCode
	}
	if newCourse.NewCourseName == "" {
		newCourse.NewCourseName = currentCourse.CourseName
	}

	existQuery = "SELECT CourseCode FROM Courses WHERE CourseCode =?"
	exists = rowExists(existQuery, newCourse.NewCourseCode)

	if exists {
		responseString = fmt.Sprintf("New course code (%s) or course name (%s) already exists.\nUnable to edit the course information.\nData has not been updated.\n", newCourse.NewCourseCode, newCourse.NewCourseName)
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(responseString))
	} else {
		editQuery := "UPDATE goms_db.Courses SET CourseCode=?, CourseName=? WHERE ID=?"
		results, err := db.Exec(editQuery, newCourse.NewCourseCode, newCourse.NewCourseName, currentCourse.ID)
		if err != nil {
			panic(err)
		} else {
			rows, _ := results.RowsAffected()
			fmt.Printf("Rows affected: %v\n", rows)
			fmt.Printf("Course: %s - %s has been modified in the database.\nNew Course code: %s - %s.\n", currentCourse.CourseCode, currentCourse.CourseName, newCourse.NewCourseCode, newCourse.NewCourseName)
			responseString = fmt.Sprintf("%s: %s has been modified to %s: %s.\n", currentCourse.CourseCode, currentCourse.CourseName, newCourse.NewCourseCode, newCourse.NewCourseName)
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(responseString))
		}
	}
}
