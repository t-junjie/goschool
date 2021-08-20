package routing

import (
	"net/http"
)

//admin is a handler that renders the admin page.
func admin(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	//redirect to appropriate menu
	if isLoggedIn(req) && !isAdmin(res, req) {
		InfoLogger.Printf("Redirecting non-admin user, %s, to user page", myUser.Username)
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "admin.gohtml", r)
	if err != nil {
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
