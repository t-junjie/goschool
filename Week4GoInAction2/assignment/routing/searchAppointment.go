package routing

import (
	"html/template"
	"net/http"
	"strings"

	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"
)

//searchAppointment is a handler that loads the appointment searching page.
//Users are redirected to login if they are not logged in.
//Users are redirected to the appointment editing page
//if they key in a valid patient name into the form.
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

		if !validate.IsAlphabet(patName) {
			WarnLogger.Println("Invalid input for patient's name. " + validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid input for patient's name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)

		}
		patient, searchErr := patients.SearchPatient(patName)
		if searchErr != nil {
			InfoLogger.Printf(searchErr.Error())
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
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
