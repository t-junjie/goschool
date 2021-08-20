package routing

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
)

func viewPatientAppointmentHistory(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser

	if isLoggedIn(req) && !isAdmin(res, req) {
		patName := myUser.FirstName + " " + myUser.LastName
		patient, searchErr := patients.SearchPatient(patName)
		if searchErr != nil {
			http.Error(res, errors.Wrap(searchErr, "No appointments made yet").Error(), http.StatusNotFound)
			return
		}
		patHist, printErr := patient.Info.AppointmentHistory.PrintAllNodes()
		if printErr != nil {
			http.Error(res, printErr.Error(), http.StatusNotFound)
		}
		r.PatientHistory = template.HTML(patHist)
	}

	err := tpl.ExecuteTemplate(res, "viewPatientAppointmentHistory.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
