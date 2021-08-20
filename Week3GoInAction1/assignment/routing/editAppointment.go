package routing

import (
	"net/http"
	"strings"
	"sync"

	helper "Assignment3/bookingSystem/helperFunctions"
)

var wgEdit sync.WaitGroup

func editAppointment(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	r.User = myUser
	r.BookingStatus = false

	if !isLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		fromDate := req.FormValue("fromDate")
		fromTime := req.FormValue("fromTime")

		fromDoctor := req.FormValue("fromDoctor")
		fromDoctor = strings.ToUpper(fromDoctor)

		toDate := req.FormValue("toDate")
		toTime := req.FormValue("toTime")

		toDoctor := req.FormValue("toDoctor")
		toDoctor = strings.ToUpper(toDoctor)

		if fromDate == "" || fromTime == "" || fromDoctor == "" || toDate == "" || toTime == "" || toDoctor == "" {
			http.Error(res, helper.ErrEmptyFields.Error(), http.StatusBadRequest)
			return
		}

		fromDT, fromErr := helper.ConvertDateTime(fromDate, fromTime)
		if fromErr != nil {
			http.Error(res, fromErr.Error(), http.StatusBadRequest)
			return
		}

		toDT, toErr := helper.ConvertDateTime(toDate, toTime)
		if toErr != nil {
			http.Error(res, toErr.Error(), http.StatusBadRequest)
			return
		}

		fDoctor, searchErr := doctors.SearchDoctor(fromDoctor)
		if searchErr != nil {
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		tDoctor, searchErr := doctors.SearchDoctor(toDoctor)
		if searchErr != nil {
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		pat, err := patients.SearchPatient(string(r.PatientName))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
		if tDoctor.IsAvailableAt(toDT) {
			wgEdit.Add(1)
			go func() {
				defer wgEdit.Done()
				tDoctor.Info.Appointments.AddAppointment(toDT, pat.Info.PatientName, tDoctor.Info.DoctorName)
				tDoctor.RemoveAvailability(toDT)
				fDoctor.Info.Appointments.RemoveAppointment(fromDT, pat.Info.PatientName)
				fDoctor.AddAvailability(fromDT)

				pat.Info.AppointmentHistory.AddAppointment(toDT, pat.Info.PatientName, tDoctor.Info.DoctorName)
				pat.Info.AppointmentHistory.RemoveAppointment(fromDT, pat.Info.PatientName)
				r.BookingStatus = true
			}()
			wgEdit.Wait()
		} else {
			http.Error(res, helper.ErrAlreadyTaken.Error(), http.StatusNotFound)
			return
		}
	}

	err := tpl.ExecuteTemplate(res, "editAppointment.gohtml", r)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
