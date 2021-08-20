package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func signup(res http.ResponseWriter, req *http.Request) {

	if isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var myUser userInfo

	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")

		firstname := req.FormValue("firstname")
		firstname = strings.ToUpper(firstname)

		lastname := req.FormValue("lastname")
		lastname = strings.ToUpper(lastname)

		if username == "" || password == "" || firstname == "" || lastname == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		if username != "" {
			if _, ok := mapUsers[username]; ok {
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}
			// create session
			id := uuid.NewV4()
			loginCookie := &http.Cookie{
				Name:  "loginCookie",
				Value: id.String(),
			}
			http.SetCookie(res, loginCookie)
			mapSessions[loginCookie.Value] = username

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = userInfo{username, bPassword, firstname, lastname}
			mapUsers[username] = myUser
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	err := tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
