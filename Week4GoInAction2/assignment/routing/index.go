package routing

import (
	"net/http"
)

//index redirects users to the login page.
//The default landing page for the application is the login page.
//If the user is logged in, the user will be redirected to either the admin page or the home page.
func index(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
}
