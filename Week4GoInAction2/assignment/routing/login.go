package routing

import (
	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//login is a handler that renders the login page.
//If the user is logged in, the user will be redirected to either the admin page or the home page.
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

		//validate inputs via a whitelist
		if !validate.IsAlphaNumeric(username) || !validate.IsAlphaNumeric(password) {
			WarnLogger.Println("Invalid username and/or password. ", validate.ErrInvalidAlphaNum.Error())
			http.Error(res, "Invalid username and/or password. "+validate.ErrInvalidAlphaNum.Error(), http.StatusUnprocessableEntity)
			return
		}

		myUser, ok := mapUsers[username]
		if !ok {
			//user does not exist
			http.Error(res, "Invalid username and/or password", http.StatusUnauthorized)
			return
		}
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			//incorrect password
			http.Error(res, "Invalid username and/or password", http.StatusForbidden)
			return
		}
		//disallow concurrent login based on active sessions in mapSessions
		for _, user := range mapSessions {
			if user == username {
				http.Error(res, "Unable to log in. Please try again later.", http.StatusSeeOther)
				return
			}
		}

		id := uuid.NewV4()
		var expireCookie = time.Now().Add(1 * time.Hour)
		loginCookie := &http.Cookie{
			Name:     "loginCookie",
			Value:    id.String(),
			Expires:  expireCookie,
			Domain:   "localhost",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
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
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
