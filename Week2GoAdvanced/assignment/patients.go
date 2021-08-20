package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type PatientInfo struct {
	patientName        string
	appointmentHistory Queue
}

type Patient struct {
	info  PatientInfo
	left  *Patient
	right *Patient
}

type Patients struct {
	root *Patient
}

var (
	ErrPatientNotFound = errors.New("The named patient was not found!")
	patientLock        sync.Mutex
)

/*----------------------------Insert Operations--------------------------*/

//insert;
//add individual node to bst
func (p *Patients) addPatient(pat **Patient, patientName, doctorName string, apptdt time.Time) error {
	if *pat == nil {
		newPatient := PatientInfo{patientName, Queue{nil, nil, 0}}
		err := newPatient.appointmentHistory.addAppointment(apptdt, patientName, doctorName)
		if err != nil {
			return errors.Wrap(err, "adding appointment failed")
		}
		newNode := &Patient{newPatient, nil, nil}
		*pat = newNode
		return nil
	}
	if patientName < (*pat).info.patientName {
		err := p.addPatient(&((*pat).left), patientName, doctorName, apptdt)
		if err != nil {
			return err
		}
	} else {
		err := p.addPatient(&((*pat).right), patientName, doctorName, apptdt)
		if err != nil {
			return err
		}
	}
	return nil
}

//wrapper function for addPatient
func (p *Patients) add(patientName, doctorName string, apptdt time.Time) error {
	patientLock.Lock()
	defer patientLock.Unlock()

	if patientName == "" {
		return ErrInvalidName
	} else {
		err := p.addPatient(&p.root, patientName, doctorName, apptdt)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

/*----------------------------Search Operations--------------------------*/

//wrapper function for searchPatNodes
func (p *Patients) searchPatient(name string) (*Patient, error) {
	return p.searchPatNodes(p.root, name)
}

//search;
//binary search in bst for node with matching name
func (p *Patients) searchPatNodes(pat *Patient, name string) (*Patient, error) {
	if pat == nil {
		return &Patient{}, ErrPatientNotFound
	} else {
		if pat.info.patientName == name {
			return pat, nil
		} else {
			if name < pat.info.patientName {
				return p.searchPatNodes(pat.left, name)
			} else {
				return p.searchPatNodes(pat.right, name)
			}
		}
	}
}

/*----------------------------Print Operations--------------------------*/

//print;
//in order traversal for bst which prints out patient's name
func (p *Patients) printPatients(pat *Patient) {
	if pat != nil {
		p.printPatients(pat.left)
		fmt.Println(pat.info.patientName)
		p.printPatients(pat.right)
	}
}

//wrapper for printPatients
func (p *Patients) print() {
	p.printPatients(p.root)
}
