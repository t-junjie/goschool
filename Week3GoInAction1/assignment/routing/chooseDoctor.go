package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"html/template"
	"net/http"
	"strings"
)

func chooseDoctor(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	tmp := doctors.Print()
	r.DoctorsName = template.HTML(tmp)
	r.DoctorName = template.HTML("")
	r.DoctorAppointments = template.HTML("")

	if isLoggedIn(req) && !isAdmin(res, req) {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		dName := req.FormValue("docname")
		dName = strings.ToUpper(dName)
		if dName == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		doctor, err := doctors.SearchDoctor(dName)
		r.DoctorName = template.HTML(doctor.Info.DoctorName)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		} else {
			docAppt, queueErr := doctor.Info.Appointments.PrintAllNodes()
			if queueErr != nil {
				r.DoctorAppointments = template.HTML("")
			}
			r.DoctorAppointments = template.HTML(docAppt)
			http.Redirect(res, req, "/viewDoctorAppointments", http.StatusSeeOther)
			return
		}
	}

	err := tpl.ExecuteTemplate(res, "chooseDoctor.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
