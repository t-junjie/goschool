package routing

import (
	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//signup is a handler that renders the signup page.
//signup redirects a user to the index page if the user is already logged in.
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

		//check for empty values
		if username == "" || password == "" || firstname == "" || lastname == "" {
			WarnLogger.Println(helper.ErrEmptyFields.Error())
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		//validate inputs via a whitelist
		if !validate.IsAlphaNumeric(username) || !validate.IsAlphaNumeric(password) {
			WarnLogger.Println("Invalid username and/or password. ", validate.ErrInvalidAlphaNum.Error())
			http.Error(res, "Invalid username and/or password. "+validate.ErrInvalidAlphaNum.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !validate.IsAlphabet(firstname) || !validate.IsAlphabet(lastname) {
			WarnLogger.Println("Invalid first/last name. ", validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid first/last name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)
		}

		if username != "" {
			if _, ok := mapUsers[username]; ok {
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}
			// create session
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
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
