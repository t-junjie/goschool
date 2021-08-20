package bookingSystem

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrEmptyQueue                          = errors.New("Appointments List is empty. There are no appointments available.")
	appointmentLock, updateAppointmentLock sync.Mutex
)

type Appointment struct {
	dateTime    time.Time
	patientName string
	doctorName  string
}

type QueueNode struct {
	appointment Appointment
	next        *QueueNode
}

type Queue struct {
	front *QueueNode
	back  *QueueNode
	size  int
}

func (q *Queue) isEmpty() bool {
	return q.size == 0
}

/*----------------------------Insert Operations--------------------------*/

//enqueue
func (q *Queue) AddAppointment(dt time.Time, patientName, doctorName string) error {
	appointmentLock.Lock()
	defer appointmentLock.Unlock()

	if dt.Hour() < workStartHour || dt.Hour() >= workEndHour {
		return errors.Wrap(ErrInvalidAppt, "outside of working hours: ")
	}
	newAppointment := Appointment{dt, patientName, doctorName}
	newNode := &QueueNode{newAppointment, nil}
	if q.front == nil {
		q.front = newNode
	} else {
		if q.front.appointment.dateTime.After(dt) { //move earlier appointment up
			newNode.next = q.front
			q.front = newNode
		} else {
			currentNode := q.front
			for currentNode.next != nil && (currentNode.next.appointment.dateTime.Before(dt) || currentNode.next.appointment.dateTime.Equal(dt)) {
				currentNode = currentNode.next
			}
			newNode.next = currentNode.next
			currentNode.next = newNode
		}
	}
	q.size++
	return nil
}

/*----------------------------Delete Operations--------------------------*/

//dequeue
func (q *Queue) FinishAppointment() (Appointment, error) {
	var appt Appointment

	if q.front == nil {
		appt = Appointment{time.Time{}, "", ""}
		return appt, ErrEmptyQueue
	}
	appt = q.front.appointment
	if q.size == 1 {
		q.front = nil
		q.back = nil
	} else {
		q.front = q.front.next
	}
	q.size--
	return appt, nil
}

//traverse a linked list to look for specified appointment and removes it
func (q *Queue) RemoveAppointment(dt time.Time, patientName string) error {
	appointmentLock.Lock()
	defer appointmentLock.Unlock()
	if q.size == 0 {
		return ErrEmptyQueue
	}
	currentNode := q.front
	time, pname := currentNode.appointment.dateTime, currentNode.appointment.patientName
	if dt == time && patientName == pname {
		_, err := q.FinishAppointment()
		if err != nil {
			return err
		}
		return nil
	}
	for currentNode.next != nil {
		time, pname = currentNode.next.appointment.dateTime, currentNode.next.appointment.patientName
		if dt == time && patientName == pname {
			currentNode.next = currentNode.next.next
			return nil
		}
		currentNode = currentNode.next
	}
	q.size--
	return errors.Wrap(ErrInvalidAppt, "did not find appointment")
}

/*----------------------------Modify Operations--------------------------*/

func (q *Queue) updateAppointment(oldappt, newappt Appointment) error {
	updateAppointmentLock.Lock()
	defer updateAppointmentLock.Unlock()
	if q.size == 0 {
		return ErrEmptyQueue
	}
	_, err := q.SearchAppointment(oldappt.dateTime, oldappt.patientName)
	if err != nil {
		return err
	} else {
		removeErr := q.RemoveAppointment(oldappt.dateTime, oldappt.patientName)
		if removeErr != nil {
			return removeErr
		}
		addErr := q.AddAppointment(newappt.dateTime, newappt.patientName, newappt.doctorName)
		if addErr != nil {
			return addErr
		}
		return nil
	}
}

/*----------------------------Search Operations--------------------------*/

//traverse a linked list to look for specified appointment
func (q *Queue) SearchAppointment(dt time.Time, patientName string) (Appointment, error) {
	var appt Appointment

	if q.size == 0 {
		appt = Appointment{time.Time{}, "", ""}
		return appt, ErrEmptyQueue
	}
	currentNode := q.front
	tm, pname, dname := currentNode.appointment.dateTime, currentNode.appointment.patientName, currentNode.appointment.doctorName
	if dt == tm && patientName == pname {
		appt = Appointment{tm, pname, dname}
		return appt, nil
	}
	for currentNode.next != nil {
		time, pname, dname := currentNode.next.appointment.dateTime, currentNode.next.appointment.patientName, currentNode.next.appointment.doctorName
		if dt == time && patientName == pname {
			appt = Appointment{time, pname, dname}
			return appt, nil
		}
		currentNode = currentNode.next
	}
	appt = Appointment{time.Time{}, "", ""}
	return appt, errors.Wrap(ErrInvalidAppt, "appointment not found")
}

/*----------------------------Print Operations--------------------------*/

//traverse a linked list and prints out appointment information
func (q *Queue) PrintAllNodes() (string, error) {
	var apptMsg string
	currentNode := q.front
	if currentNode == nil {
		return "", ErrEmptyQueue
	}

	apptMsg = fmt.Sprintf("%s's appointment with DR. %s is at %v\n<br>", strings.ToUpper(currentNode.appointment.patientName), currentNode.appointment.doctorName, currentNode.appointment.dateTime.Format("Mon, 02 Jan 2006 15:04:05"))

	for currentNode.next != nil {
		currentNode = currentNode.next
		apptMsg = apptMsg + fmt.Sprintf("%s's appointment with DR. %s is at %v\n<br>", strings.ToUpper(currentNode.appointment.patientName), currentNode.appointment.doctorName, currentNode.appointment.dateTime.Format("Mon, 02 Jan 2006 15:04:05"))
	}
	return apptMsg, nil
}

func (q *QueueNode) Print() string {
	return fmt.Sprintf("%s's appointment with DR. %s is at %v\n<br>", strings.ToUpper(q.appointment.patientName), strings.ToUpper(q.appointment.doctorName), q.appointment.dateTime)
}
