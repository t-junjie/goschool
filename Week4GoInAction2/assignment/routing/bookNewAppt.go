package routing

import (
	"net/http"
	"strings"
	"sync"

	helper "Assignment4/booking/helperFunctions"

	validate "Assignment4/validation"
)

var wg sync.WaitGroup

//bookNewAppointment is a handler that renders the booking page.
func bookNewAppointment(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	r.BookingStatus = false

	//redirect if not user
	if isLoggedIn(req) && isAdmin(res, req) {
		InfoLogger.Printf("Redirecting %s to admin page", myUser.Username)
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	}

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		var patientExists bool

		date := req.FormValue("date")
		time := req.FormValue("time")

		docName := req.FormValue("name")
		docName = strings.ToUpper(docName)

		if date == "" || time == "" || docName == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		if !validate.IsDateString(date) {
			WarnLogger.Println(validate.ErrInvalidDate.Error())
			http.Error(res, validate.ErrInvalidDate.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !validate.IsTimeString(time) {
			WarnLogger.Println(validate.ErrInvalidTime.Error())
			http.Error(res, validate.ErrInvalidTime.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !validate.IsAlphabet(docName) {
			WarnLogger.Println("Invalid input for doctor's name. " + validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid input for doctor's name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)
		}

		doctor, searchErr := doctors.SearchDoctor(docName)
		if searchErr != nil {
			InfoLogger.Printf("DR. %s - %s", docName, searchErr.Error())
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		patientName := r.User.FirstName + " " + r.User.LastName

		apptDate, convertErr := helper.ConvertDateTime(date, time)
		if convertErr != nil {
			InfoLogger.Println(convertErr.Error())
			http.Error(res, convertErr.Error(), http.StatusBadRequest)
			return
		}

		patient, searchErr := patients.SearchPatient(patientName)
		if searchErr != nil { //patient is not found
			InfoLogger.Println(searchErr.Error())
			patientExists = false
		} else {
			patientExists = true
		}

		if doctor.IsAvailableAt(apptDate) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				doctor.Info.Appointments.AddAppointment(apptDate, patientName, doctor.Info.DoctorName)
				doctor.RemoveAvailability(apptDate)
				if !patientExists {
					patients.Add(patientName, doctor.Info.DoctorName, apptDate)
				} else {
					patient.Info.AppointmentHistory.AddAppointment(apptDate, patientName, doctor.Info.DoctorName)
				}
				r.BookingStatus = true
			}()
			wg.Wait()
		} else {
			http.Error(res, helper.ErrAlreadyTaken.Error(), http.StatusNotFound)
			return
		}
	}
	err := tpl.ExecuteTemplate(res, "booking.gohtml", r)
	if err != nil {
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
