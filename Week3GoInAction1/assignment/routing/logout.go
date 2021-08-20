package routing

import "net/http"

func logout(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	loginCookie, _ := req.Cookie("loginCookie")
	// delete the session
	delete(mapSessions, loginCookie.Value)
	// remove the cookie
	loginCookie = &http.Cookie{
		Name:   "loginCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, loginCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}
