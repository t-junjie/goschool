package routing

import (
	helper "Assignment3/bookingSystem/helperFunctions"
	"html/template"
	"net/http"
	"strings"
)

func choosePatient(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	tmp := patients.Print()
	r.PatientsName = template.HTML(tmp)
	r.PatientHistory = template.HTML("")

	if isLoggedIn(req) && !isAdmin(res, req) {
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		pName := req.FormValue("patname")
		pName = strings.ToUpper(pName)
		if pName == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		patient, err := patients.SearchPatient(pName)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		} else {
			patHist, queueErr := patient.Info.AppointmentHistory.PrintAllNodes()
			if queueErr != nil {
				http.Error(res, queueErr.Error(), http.StatusNotFound)
				return
			}
			r.PatientHistory = template.HTML(patHist)
			http.Redirect(res, req, "/viewPatientAppointmentHistory", http.StatusSeeOther)
			return
		}
	}

	err := tpl.ExecuteTemplate(res, "choosePatient.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
