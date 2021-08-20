package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"html/template"
	"net/http"
)

func searchAvailableDoctors(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	//reinitialize availability results to default values
	r.DoctorsAvailability = template.HTML("")
	r.DoctorAvailability = template.HTML("")

	if isLoggedIn(req) && isAdmin(res, req) {
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		dt := req.FormValue("datetime")
		if dt == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		date, convertErr := helper.ConvertDateTime(dt, "0000")
		if convertErr != nil {
			http.Error(res, convertErr.Error(), http.StatusBadRequest)
			return
		}
		docsAvail := doctors.ShowAllAvailability(date)
		r.DoctorsAvailability = template.HTML(docsAvail)
		http.Redirect(res, req, "/searchDoctorAvailability", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "searchAvailableDoctors.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
