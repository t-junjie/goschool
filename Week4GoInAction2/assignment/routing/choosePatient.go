package routing

import (
	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"
	"html/template"
	"net/http"
	"strings"
)

//choosePatient is a handler that renders the patient selection page.
func choosePatient(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	tmp := patients.Print()
	r.PatientsName = template.HTML(tmp)
	r.PatientHistory = template.HTML("")

	if isLoggedIn(req) && !isAdmin(res, req) {
		InfoLogger.Printf("Redirecting non-admin user, %s, to user page", myUser.Username)
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

		if !validate.IsAlphabet(pName) {
			WarnLogger.Println("Invalid input for patient's name. " + validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid input for patient's name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)
			return
		}

		patient, err := patients.SearchPatient(pName)
		if err != nil {
			InfoLogger.Printf(err.Error())
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
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
