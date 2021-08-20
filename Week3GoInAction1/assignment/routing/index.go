package routing

import (
	"net/http"
)

func index(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
}
