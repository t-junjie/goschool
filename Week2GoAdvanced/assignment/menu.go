package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidDate     = errors.New("Invalid Date Entered")
	ErrInvalidOption   = errors.New("Invalid option chosen. Please retry.")
	ErrInsufficientArg = errors.New("Missing information. Please check your input and retry.")
	ErrAlreadyTaken    = errors.New("This time slot is unavailable. Please retry.")
	wg                 sync.WaitGroup
)

const (
	userPrompt string = `

Select one of the following options:
1. Book a new appointment
2. Change existing appointment
3. Search for available doctors 
4. List available times of selected doctor
`

	adminPrompt string = `

Select one of the following options:
1. Browse appointments for a doctor
2. View patient's appointment history
`

	datePrompt string = `

Please enter the date to check the doctors' availability
*Only weekdays from next week onwards to the week after that is available.

Acceptable Format: YYYYMMDD e.g. 20210417 
(If no doctors are available, nothing will be displayed.)
`
	getDocNamePrompt string = `

Please enter the name of the doctor you wish to book an appointment with (case insensitive).
e.g. LUO, GLENN 
`

	browseDocApptPrompt string = `

Which doctor's schedule would you like to view?
`
	browsePatientApptPrompt string = `

Which patient's appointment history would you like to view?
`
	bookingPrompt string = `

Please enter your name, the date and time you would like the appointment to be at (case insensitive).
e.g. Jerome Chew Boon Kah; 20210418 (YYYYMMDD); 1530 (HHMM)
`
	getBookingPrompt string = `
Please enter the name of the doctor you previously booked an appointment with (case insensitive).
e.g. LUO, GLENN  
`
	getPastBookingPrompt string = `

Please enter your name, the date and time you have booked the appointment to be at (case insensitive).
e.g. Jerome Chew Boon Kah; 20210418 (YYYYMMDD); 1530 (HHMM)
`
)

func runAdminMenu(d *Doctors, p *Patients) {
	var adminInput int
	fmt.Scanln(&adminInput)

	switch adminInput {
	case 1:
		fmt.Println(browseDocApptPrompt)
		d.print()
		fmt.Println()

		name := strings.ToUpper(getInput())
		doctor, err := d.searchDoctor(name)
		if err != nil {
			fmt.Println(err)
		} else {
			doctor.info.appointments.printAllNodes()
		}

	case 2:
		fmt.Println(browsePatientApptPrompt)
		p.print()
		fmt.Println()

		name := strings.ToUpper(getInput())
		patient, err := p.searchPatient(name)
		if err != nil {
			fmt.Println(err)
		} else {
			patient.info.appointmentHistory.printAllNodes()
		}
	default:
		fmt.Println("Invalid option. Returning to main menu.")
	}
	time.Sleep(1000 * time.Millisecond)
}

func runUserMenu(d *Doctors, p *Patients) {
	var userInput int
	fmt.Scanln(&userInput)

	switch userInput {
	case 1:
		fmt.Println("Make a new appointment")
		err := makeNewBooking(d, p)
		if err != nil {
			fmt.Println(err)
		}
	case 2:
		fmt.Println("Edit existing appointment")
		err := editBooking(d, p)
		if err != nil {
			fmt.Println(err)
		}
	case 3:
		fmt.Println(datePrompt)
		var dt string
		fmt.Scanln(&dt)
		for len(dt) != 8 {
			fmt.Println(datePrompt)
			fmt.Scanln(&dt)
		}
		yy, _ := strconv.Atoi(dt[:4])
		mm, _ := strconv.Atoi(dt[4:6])
		dd, _ := strconv.Atoi(dt[6:])
		tmpDate := time.Date(yy, time.Month(mm), dd, 0, 0, 0, 0, time.Local)
		d.showAllAvailability(tmpDate)
	case 4:
		fmt.Println(getDocNamePrompt)
		name := strings.ToUpper(getInput())
		doctor, err := d.searchDoctor(name)
		if err != nil {
			fmt.Println(err)
		} else {
			doctor.showAvailability()
		}
	default:
		fmt.Println("Invalid option. Returning to main menu.")
	}
	time.Sleep(1000 * time.Millisecond)
}

func makeNewBooking(d *Doctors, p *Patients) error {
	var patientExists bool
	fmt.Println(getDocNamePrompt)
	name := strings.ToUpper(getInput())
	doctor, searchErr := d.searchDoctor(name)
	if searchErr != nil {
		return searchErr
	}

	patientName, date, time, bookingErr := getBookingDetails(bookingPrompt)
	if bookingErr != nil {
		return bookingErr
	}

	apptDate, convertErr := convertDateTime(date, time)
	if convertErr != nil {
		return convertErr
	}

	patient, searchErr := p.searchPatient(patientName)
	if searchErr != nil { //patient is not found
		patientExists = false
	} else {
		patientExists = true
	}

	if doctor.isAvailableAt(apptDate) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doctor.info.appointments.addAppointment(apptDate, patientName, doctor.info.doctorName)
			doctor.removeAvailability(apptDate)
			if !patientExists {
				p.add(patientName, doctor.info.doctorName, apptDate)
			} else {
				patient.info.appointmentHistory.addAppointment(apptDate, patientName, doctor.info.doctorName)
			}
			fmt.Println("Appointment Booked")
		}()
		wg.Wait()
	} else {
		return ErrAlreadyTaken
	}
	return nil
}

//assume rebook to change time
func editBooking(d *Doctors, p *Patients) error {
	fmt.Println(getBookingPrompt)
	name := strings.ToUpper(getInput())
	doctor, searchErr := d.searchDoctor(name)
	if searchErr != nil {
		return searchErr
	}

	patientName, date, time, bookingErr := getBookingDetails(getPastBookingPrompt)
	if bookingErr != nil {
		return bookingErr
	}

	apptDate, convertErr := convertDateTime(date, time)
	if convertErr != nil {
		return convertErr
	}
	//remove past booking
	oldAppt, err := doctor.info.appointments.searchAppointment(apptDate, patientName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Requested Appointment Date - %s, Patient Name - %s", apptDate.Format("Mon, 02 Jan 2006 15:04:05"), patientName))
	}

	patientName, date, time, bookingErr = getBookingDetails(bookingPrompt)
	if bookingErr != nil {
		return bookingErr
	}

	apptDate, convertErr = convertDateTime(date, time)
	if convertErr != nil {
		return convertErr
	}

	newAppt := Appointment{apptDate, patientName, doctor.info.doctorName}
	patient, searchErr := p.searchPatient(patientName)
	if searchErr != nil {
		return searchErr
	}

	if doctor.isAvailableAt(newAppt.dateTime) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doctor.info.appointments.updateAppointment(oldAppt, newAppt)
			doctor.addAvailability(oldAppt.dateTime)
			doctor.removeAvailability(newAppt.dateTime)
			patient.info.appointmentHistory.updateAppointment(oldAppt, newAppt)
			fmt.Println("Appointment changed")
		}()
		wg.Wait()
	} else {
		return ErrAlreadyTaken
	}

	return nil

}

func convertDateTime(date, tm string) (time.Time, error) {
	if len(date) != 8 || len(tm) != 4 {
		return time.Time{}, errors.Wrap(ErrInvalidDate, "incorrect date or time")
	}
	yy, _ := strconv.Atoi(date[:4])
	mm, _ := strconv.Atoi(date[4:6])
	dd, _ := strconv.Atoi(date[6:])
	hh, _ := strconv.Atoi(tm[:2])
	m, _ := strconv.Atoi(tm[2:])
	return time.Date(yy, time.Month(mm), dd, hh, m, 0, 0, time.Local), nil
}

func getBookingDetails(prompt string) (patientName, date, tm string, e error) {

	var input string
	fmt.Println(prompt)
	input = getInput()
	inputFields := strings.Split(input, ";")
	if len(inputFields) != 3 {
		return "", "", "", ErrInsufficientArg
	}
	patientName = inputFields[0]
	date = strings.ReplaceAll(inputFields[1], " ", "")
	tm = strings.ReplaceAll(inputFields[2], " ", "")
	return strings.ToUpper(patientName), date, tm, nil
}
