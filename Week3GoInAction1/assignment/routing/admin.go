package routing

import (
	"net/http"
)

func admin(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	//redirect to appropriate menu
	if isLoggedIn(req) && !isAdmin(res, req) {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "admin.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
