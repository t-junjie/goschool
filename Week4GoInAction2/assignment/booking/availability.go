package booking

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//These define a doctor's default availability.
//Consultation hours are from 09:00-16:00.
//No consultation available during lunch time, 12:00.
//Available dates start from one week from current date to the week after.
const (
	workStartHour   int = 9
	workEndHour     int = 16
	oneWeekFromNow  int = 7
	twoWeeksFromNow int = 15
	lunchHour       int = 12
)

//initAvailability returns a default slice of available time slots.
func initAvailability() []time.Time {
	yy, mm, dd := time.Now().Date()
	//endOfMonth := time.Date(yy, mm+1, 0, workEndHour, 0, 0, 0, time.Now().Local().Location())
	timeSlots := make([]time.Time, 0)
	for i := oneWeekFromNow; i < twoWeeksFromNow; i++ { //create Availability for the next week
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

//AddAvailability frees up a time slot in a doctor's schedule.
func (d *Doctor) AddAvailability(dt time.Time) error {
	availabilityLock.Lock()
	defer availabilityLock.Unlock()
	hour, _, _ := dt.Clock()
	if hour < workStartHour || hour >= workEndHour {
		return errors.Wrap(ErrInvalidAppt, "outside of working hours: ")
	} else {
		d.Info.Availability = Insert(d.Info.Availability, dt)
		return nil
	}
}

//Insert adds a time slot to a slice of time slots in sorted order.
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

//RemoveAvailability removes a time slot in a doctor's schedule.
func (d *Doctor) RemoveAvailability(dt time.Time) error {
	availabilityLock.Lock()
	defer availabilityLock.Unlock()
	for i, timeslot := range d.Info.Availability {
		if timeslot == dt {
			d.Info.Availability = append(d.Info.Availability[:i], d.Info.Availability[i+1:]...)
			return nil
		}
	}
	return errors.Wrap(ErrInvalidAppt, "no corresponding appointment found")
}

/*----------------------------Print Operations--------------------------*/

//showAvailabilityOn recursively checks if a doctor is available on a particular date.
func (d *Doctors) showAvailabilityOn(doc *Doctor, dt time.Time) string {
	if doc != nil {
		//d.showAvailabilityOn(doc.left, dt)
		if doc.IsAvailableOn(dt) {
			return d.showAvailabilityOn(doc.left, dt) + fmt.Sprintf("DR.%s is available on %v\n", doc.Info.DoctorName, dt.Format("Mon, 02 Jan 2006")) + "<br>" + d.showAvailabilityOn(doc.right, dt)
		}
		//d.showAvailabilityOn(doc.right, dt)
	}
	return ""
}

//ShowAllAvailability returns a string of doctors that are available on a particular date.
func (d *Doctors) ShowAllAvailability(dt time.Time) string {
	return d.showAvailabilityOn(d.root, dt)
}

//ShowAvailability shows all timeslots for which a doctor is available.
func (d *Doctor) ShowAvailability() string {
	var avail []string
	if d == nil {
		fmt.Print(ErrDocNotFound.Error())
	} else {
		noOfTimeSlots := len(d.Info.Availability)
		for _, dates := range d.Info.Availability {
			dt := dates.Format("Mon, 02 Jan 2006 15:04:05")
			avail = append(avail, dt)
		}
		return fmt.Sprintf("<h4>DR.%s is available at the following %d timeslots: <br></h4>%s", d.Info.DoctorName, noOfTimeSlots, strings.Join(avail, "<br>"))
	}
	return ""
}

/*----------------------------Helper Operations--------------------------*/

//IsAvailableOn shows the availability of a doctor on a particular day.
func (d *Doctor) IsAvailableOn(dt time.Time) bool {
	//availble on a particular date
	for _, date := range d.Info.Availability {
		requiredYY, requiredMM, requiredDD := dt.Date()
		givenYY, givenMM, givenDD := date.Date()
		if requiredYY == givenYY && requiredMM == givenMM && requiredDD == givenDD {
			return true
		}
	}
	return false
}

//IsAvailableAt shows the availability of a doctor on a particular date and time.
func (d *Doctor) IsAvailableAt(dt time.Time) bool {
	//availble on a particular date
	for _, date := range d.Info.Availability {
		if dt.Equal(date) {
			return true
		}
	}
	return false
}
