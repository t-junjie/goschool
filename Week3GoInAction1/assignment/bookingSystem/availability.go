package bookingSystem

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
	AvailabilityLock sync.Mutex
)

func initAvailability() []time.Time {
	//initialize doctor Availability from nextWeek to nextWeek + 1 day (to make viewing Availability easier)
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

func (d *Doctor) AddAvailability(dt time.Time) error {
	AvailabilityLock.Lock()
	defer AvailabilityLock.Unlock()
	hour, _, _ := dt.Clock()
	if hour < workStartHour || hour >= workEndHour {
		return errors.Wrap(ErrInvalidAppt, "outside of working hours: ")
	} else {
		d.Info.Availability = Insert(d.Info.Availability, dt)
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

func (d *Doctor) RemoveAvailability(dt time.Time) error {
	AvailabilityLock.Lock()
	defer AvailabilityLock.Unlock()
	for i, timeslot := range d.Info.Availability {
		if timeslot == dt {
			d.Info.Availability = append(d.Info.Availability[:i], d.Info.Availability[i+1:]...)
			return nil
		}
	}
	return errors.Wrap(ErrInvalidAppt, "no corresponding appointment found")
}

/*----------------------------Print Operations--------------------------*/

//show all doctors that are available on a particular date
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

//show all doctors that are available on a particular date
func (d *Doctors) ShowAllAvailability(dt time.Time) string {
	return d.showAvailabilityOn(d.root, dt)
}

//show all timeslots for which a single doctor is available
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

//Shows Availability of a doctor n a particular Date
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

//Shows Availability of a doctor on a particular Date and Time
func (d *Doctor) IsAvailableAt(dt time.Time) bool {
	//availble on a particular date
	for _, date := range d.Info.Availability {
		if dt.Equal(date) {
			return true
		}
	}
	return false
}
