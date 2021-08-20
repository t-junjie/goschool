package routing

import (
	"net/http"
)

func viewDoctorAppointments(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser

	if isLoggedIn(req) && !isAdmin(res, req) {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "viewDoctorAppointments.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
