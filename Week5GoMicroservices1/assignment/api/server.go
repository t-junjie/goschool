// Package api is a REST API that allows users to perform CRUD operations on courses.
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbCredFileName = "dbCredentials.txt"
	dbPort         = "5001"
)

var (
	db *sql.DB //Common database used by REST API to perform CRUD actions for courses
)

type Course struct {
	ID         int
	CourseCode string
	CourseName string
}

type courseInfo struct {
	CourseName string `json:"CourseName"`
}

type newCourseInfo struct {
	NewCourseCode string `json:"NewCourseCode"`
	NewCourseName string `json:"NewCourseName"`
}

func main() {

	// read username/password from file and create a valid DSN for connecting to the database
	dbCred := readCredentials(dbCredFileName)
	s := strings.Split(dbCred, ",")
	username, password, table := s[0], s[1], strings.TrimSuffix(s[2], "\n")
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", username, password, dbPort, table)

	// Use mysql as driverName and a valid DSN as dataSourceName
	var err error
	db, err = sql.Open("mysql", dsn)
	defer db.Close()
	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Successfully connected to the database on port %s...\n", dbPort)
	}

	// set up routing for REST API
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", index)
	router.HandleFunc("/api/v1/courses", viewall)
	router.HandleFunc("/api/v1/courses/{courseCode}", course).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Starting REST API server at port 5221...")
	http.ListenAndServe(":5221", router)
}

// rowExists checks if a certain row has been created in the database
func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

// readCredentials reads a named file in the current directory.
func readCredentials(filename string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	filepath := fmt.Sprintf("%s/%s", basepath, filename)

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Failed to open DB credentials file.")
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read DB credentials file.")
		panic(err)
	}
	credentials := string(bytes)
	return credentials
}
