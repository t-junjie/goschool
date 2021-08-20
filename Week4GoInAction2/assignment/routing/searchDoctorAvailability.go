package routing

import (
	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"

	"html/template"
	"net/http"
	"strings"
)

//searchDoctorAvailability is a handler that displays the doctor's availability page.
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

		if !validate.IsAlphabet(name) {
			WarnLogger.Println("Invalid input for doctor's name. " + validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid input for doctor's name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)
			return
		}

		doc, docErr := doctors.SearchDoctor(name)
		if docErr != nil {
			InfoLogger.Printf("DR. %s - %s", name, docErr.Error())
			http.Error(res, docErr.Error(), http.StatusNotFound)
			return
		}
		docAvail := doc.ShowAvailability()
		r.DoctorAvailability = template.HTML(docAvail)
	}

	err := tpl.ExecuteTemplate(res, "searchDoctorAvailability.gohtml", r)
	if err != nil {
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
