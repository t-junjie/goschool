package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// course allows its caller to perform CRUD operations on the
// course data stored in the database.
func course(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	courseCode := params["courseCode"]

	// GET is for retrieving information on a course
	if req.Method == "GET" {
		getCourse(db, res, courseCode)
	}

	// DELETE is for removing information on a course
	if req.Method == "DELETE" {
		delCourse(db, res, courseCode)
	}

	if req.Header.Get("Content-type") == "application/json" {

		// POST is for adding a new course
		if req.Method == "POST" {
			addCourse(db, res, req, courseCode)
		}

		// PUT is for editing information on an existing course
		if req.Method == "PUT" {
			editCourse(db, res, req, courseCode)
		}
	}
}
