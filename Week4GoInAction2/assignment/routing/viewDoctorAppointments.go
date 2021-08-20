package routing

import (
	"net/http"
)

//viewDoctorAppointments is a handler that renders a doctor's schedule on the page.
func viewDoctorAppointments(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser

	if isLoggedIn(req) && !isAdmin(res, req) {
		InfoLogger.Printf("Redirecting non-admin user, %s, to user page", myUser.Username)
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "viewDoctorAppointments.gohtml", r)
	if err != nil {
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
