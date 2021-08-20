package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	workStartHour   int = 9
	workEndHour     int = 15
	oneWeekFromNow  int = 7
	twoWeeksFromNow int = 15 //reduce timespan so it's easier to read debug output
	lunchHour       int = 12
)

var (
	ErrInvalidAppt   = errors.New("Invalid Appointment")
	availabilityLock sync.Mutex
)

func initAvailability() []time.Time {
	//initialize doctor availability from nextWeek to nextWeek + 1 day (to make viewing availability easier)
	yy, mm, dd := time.Now().Date()
	//endOfMonth := time.Date(yy, mm+1, 0, workEndHour, 0, 0, 0, time.Now().Local().Location())
	timeSlots := make([]time.Time, 0)
	for i := oneWeekFromNow; i < twoWeeksFromNow; i++ { //create availability for the next week
		workStart := time.Date(yy, mm, dd+i, workStartHour, 0, 0, 0, time.Now().Local().Location())
		workEnd := time.Date(yy, mm, dd+i, workEndHour, 0, 0, 0, time.Now().Local().Location())

		for workStart.Before(workEnd) && workStart.Weekday() != time.Sunday && workStart.Weekday() != time.Saturday {
			if workStart.Hour() == lunchHour {
				workStart = workStart.Add(time.Hour)
			}
			timeSlots = append(timeSlots, workStart)
			workStart = workStart.Add(time.Minute * 30)
		}
	}
	return timeSlots
}

/*----------------------------Insert Operations--------------------------*/

func (d *Doctor) addAvailability(dt time.Time) error {
	availabilityLock.Lock()
	defer availabilityLock.Unlock()
	hour, _, _ := dt.Clock()
	if hour < workStartHour || hour >= workEndHour {
		return errors.Wrap(ErrInvalidAppt, "outside of working hours: ")
	} else {
		d.info.availability = Insert(d.info.availability, dt)
		return nil
	}
}

func Insert(arr []time.Time, t time.Time) []time.Time {
	for i := 0; i < len(arr); i++ {
		if t.After(arr[len(arr)-1]) {
			arr = append(arr, t)
			break
		} else if t.Before(arr[0]) {
			arr = append(arr, t) //insert arbitary time value
			copy(arr[1:], arr[:])
			arr[0] = t
			break
		} else if t.After(arr[i]) && t.Before(arr[i+1]) {
			arr = append(arr, t)
			copy(arr[i+1:], arr[i:])
			arr[i] = t
			break
		}
	}
	return arr
}

/*----------------------------Remove Operations--------------------------*/

func (d *Doctor) removeAvailability(dt time.Time) error {
	availabilityLock.Lock()
	defer availabilityLock.Unlock()
	for i, timeslot := range d.info.availability {
		if timeslot == dt {
			d.info.availability = append(d.info.availability[:i], d.info.availability[i+1:]...)
			return nil
		}
	}
	return errors.Wrap(ErrInvalidAppt, "no corresponding appointment found")
}

/*----------------------------Print Operations--------------------------*/

//show all doctors that are available on a particular date
func (d *Doctors) showAvailabilityOn(doc *Doctor, dt time.Time) {
	if doc != nil {
		d.showAvailabilityOn(doc.left, dt)
		if doc.isAvailableOn(dt) {
			fmt.Printf("DR.%s is available on %v\n", doc.info.doctorName, dt.Format("Mon, 02 Jan 2006"))
		}
		d.showAvailabilityOn(doc.right, dt)
	}
}

//wrapper function for showAvailabilityOn
func (d *Doctors) showAllAvailability(dt time.Time) {
	d.showAvailabilityOn(d.root, dt)
}

//show all timeslots for which a single doctor is available
func (d *Doctor) showAvailability() {
	var avail []string
	if d == nil {
		fmt.Print(ErrDocNotFound.Error())
	} else {
		noOfTimeSlots := len(d.info.availability)
		for _, dates := range d.info.availability {
			avail = append(avail, dates.Format("Mon, 02 Jan 2006 15:04:05"))
		}
		fmt.Printf("DR.%s is available at the following %d timeslots: \n%s", d.info.doctorName, noOfTimeSlots, strings.Join(avail, "\n"))
	}
}

/*----------------------------Helper Operations--------------------------*/

//Shows availability of a doctor n a particular Date
func (d *Doctor) isAvailableOn(dt time.Time) bool {
	//availble on a particular date
	for _, date := range d.info.availability {
		requiredYY, requiredMM, requiredDD := dt.Date()
		givenYY, givenMM, givenDD := date.Date()
		if requiredYY == givenYY && requiredMM == givenMM && requiredDD == givenDD {
			return true
		}
	}
	return false
}

//Shows availability of a doctor on a particular Date and Time
func (d *Doctor) isAvailableAt(dt time.Time) bool {
	//availble on a particular date
	for _, date := range d.info.availability {
		if dt.Equal(date) {
			return true
		}
	}
	return false
}
