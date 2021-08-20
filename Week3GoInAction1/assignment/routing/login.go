package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func login(res http.ResponseWriter, req *http.Request) {

	if isLoggedIn(req) && isAdmin(res, req) {
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	} else if isLoggedIn(req) {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		if username == "" || password == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		myUser, ok := mapUsers[username]
		if !ok {
			http.Error(res, "User does not exist. Please navigate back and sign up.", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		id := uuid.NewV4()
		loginCookie := &http.Cookie{
			Name:  "loginCookie",
			Value: id.String(),
		}
		http.SetCookie(res, loginCookie)
		mapSessions[loginCookie.Value] = username
		if myUser.Username == "admin" {
			http.Redirect(res, req, "/admin", http.StatusSeeOther)
			return
		} else {
			http.Redirect(res, req, "/home", http.StatusSeeOther)
			return
		}
	}
	err := tpl.ExecuteTemplate(res, "login.gohtml", nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
