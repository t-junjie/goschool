package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

//delCourse allows users to delete a course
func delCourse(db *sql.DB, res http.ResponseWriter, courseCode string) {
	var responseString string

	existQuery := "SELECT CourseCode FROM Courses WHERE CourseCode=?"
	exists := rowExists(existQuery, courseCode)

	if exists {
		query := "DELETE FROM Courses WHERE CourseCode=?"
		_, err := db.Exec(query, courseCode)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Course code: %s has been deleted from the database.\n", courseCode)

		responseString = fmt.Sprintf("%s was successfully deleted.\n", courseCode)
		res.WriteHeader(http.StatusAccepted)
		res.Write([]byte(responseString))
	} else {
		responseString = "HTTP 404: Course not found."
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(responseString))
	}
}
