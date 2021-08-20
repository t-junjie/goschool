package routing

import (
	"html/template"
	"net/http"
	"strings"

	helper "Assignment3/bookingSystem/helperFunctions"
)

func searchAppointment(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	r.PatientName = template.HTML("")
	r.PatientHistory = template.HTML("")

	tmp := patients.Print()
	r.PatientsName = template.HTML(tmp)

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		patName := req.FormValue("patname")
		patName = strings.ToUpper(patName)
		if patName == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		patient, searchErr := patients.SearchPatient(patName)
		if searchErr != nil {
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}
		name := patient.Info.PatientName
		r.PatientName = template.HTML(name)

		patHist, err := patient.Info.AppointmentHistory.PrintAllNodes()
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		r.PatientHistory = template.HTML(patHist)
		http.Redirect(res, req, "/editAppointment", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(res, "searchAppointment.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
