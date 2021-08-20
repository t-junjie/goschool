package routing

import (
	"net/http"
)

func user(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	//redirect to appropriate menu
	if isLoggedIn(req) && isAdmin(res, req) {
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	}

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "user.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
