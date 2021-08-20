package main

import (
	routing "Assignment4/routing"
)

//Application runs on this port in localhost
const port = 5221

//main calls the router to start the server.
func main() {
	routing.Route(port)
}
