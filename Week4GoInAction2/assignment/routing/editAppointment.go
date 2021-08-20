package routing

import (
	"net/http"
	"strings"
	"sync"

	helper "Assignment4/booking/helperFunctions"
	validate "Assignment4/validation"
)

var wgEdit sync.WaitGroup

//editAppointment is a handler that renders the appointment editing page.
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

		if !validate.IsDateString(fromDate) || !validate.IsDateString(toDate) {
			WarnLogger.Println(validate.ErrInvalidDate.Error())
			http.Error(res, validate.ErrInvalidDate.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !validate.IsTimeString(fromTime) || !validate.IsTimeString(toTime) {
			WarnLogger.Println(validate.ErrInvalidTime.Error())
			http.Error(res, validate.ErrInvalidTime.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !validate.IsAlphabet(toDoctor) || !validate.IsAlphabet(fromDoctor) {
			WarnLogger.Println("Invalid input for doctor's name. " + validate.ErrInvalidAlphabet.Error())
			http.Error(res, "Invalid input for doctor's name. "+validate.ErrInvalidAlphabet.Error(), http.StatusUnprocessableEntity)
			return
		}

		fromDT, fromErr := helper.ConvertDateTime(fromDate, fromTime)
		if fromErr != nil {
			InfoLogger.Println(fromErr.Error())
			http.Error(res, fromErr.Error(), http.StatusBadRequest)
			return
		}

		toDT, toErr := helper.ConvertDateTime(toDate, toTime)
		if toErr != nil {
			InfoLogger.Println(toErr.Error())
			http.Error(res, toErr.Error(), http.StatusBadRequest)
			return
		}

		fDoctor, searchErr := doctors.SearchDoctor(fromDoctor)
		if searchErr != nil {
			InfoLogger.Printf("DR. %s - %s", fromDoctor, searchErr.Error())
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		tDoctor, searchErr := doctors.SearchDoctor(toDoctor)
		if searchErr != nil {
			InfoLogger.Printf("DR. %s - %s", toDoctor, searchErr.Error())
			http.Error(res, searchErr.Error(), http.StatusNotFound)
			return
		}

		pat, err := patients.SearchPatient(string(r.PatientName))
		if err != nil {
			InfoLogger.Printf("Patient %s - %s", string(r.PatientName), err.Error())
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
		ErrLogger.Println("Failed to load template: ", err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
