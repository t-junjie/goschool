package routing

import (
	"net/http"
	"strings"
	"sync"

	helper "Assignment3/bookingSystem/helperFunctions"
)

var wg sync.WaitGroup

func bookNewAppointment(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	r.BookingStatus = false

	//redirect if not user
	if isLoggedIn(req) && isAdmin(res, req) {
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

		doctor, searchErr := doctors.SearchDoctor(docName)
		if searchErr != nil {
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		patientName := r.User.FirstName + " " + r.User.LastName

		apptDate, convertErr := helper.ConvertDateTime(date, time)
		if convertErr != nil {
			http.Error(res, convertErr.Error(), http.StatusBadRequest)
			return
		}

		patient, searchErr := patients.SearchPatient(patientName)
		if searchErr != nil { //patient is not found
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
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
