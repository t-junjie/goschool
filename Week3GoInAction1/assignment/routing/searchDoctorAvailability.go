package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"html/template"
	"net/http"
	"strings"
)

func searchDoctorAvailability(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser

	if isLoggedIn(req) && isAdmin(res, req) {
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		name := req.FormValue("docname")
		name = strings.ToUpper(name)

		if name == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		doc, docErr := doctors.SearchDoctor(name)
		if docErr != nil {
			http.Error(res, docErr.Error(), http.StatusNotFound)
		}
		docAvail := doc.ShowAvailability()
		r.DoctorAvailability = template.HTML(docAvail)
	}

	err := tpl.ExecuteTemplate(res, "searchDoctorAvailability.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
