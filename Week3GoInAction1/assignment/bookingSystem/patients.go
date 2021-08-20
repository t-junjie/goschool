package bookingSystem

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type PatientInfo struct {
	PatientName        string
	AppointmentHistory Queue
}

type Patient struct {
	Info  PatientInfo
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

/*----------------------------Helper Functions--------------------------*/

func CreatePatients() Patients {
	return Patients{nil}
}

/*----------------------------Insert Operations--------------------------*/

//insert;
//add individual node to bst
func (p *Patients) addPatient(pat **Patient, PatientName, doctorName string, apptdt time.Time) error {
	if *pat == nil {
		newPatient := PatientInfo{PatientName, Queue{nil, nil, 0}}
		err := newPatient.AppointmentHistory.AddAppointment(apptdt, PatientName, doctorName)
		if err != nil {
			return errors.Wrap(err, "adding appointment failed")
		}
		newNode := &Patient{newPatient, nil, nil}
		*pat = newNode
		return nil
	}
	if PatientName < (*pat).Info.PatientName {
		err := p.addPatient(&((*pat).left), PatientName, doctorName, apptdt)
		if err != nil {
			return err
		}
	} else {
		err := p.addPatient(&((*pat).right), PatientName, doctorName, apptdt)
		if err != nil {
			return err
		}
	}
	return nil
}

//wrapper function for addPatient
func (p *Patients) Add(PatientName, doctorName string, apptdt time.Time) error {
	patientLock.Lock()
	defer patientLock.Unlock()

	if PatientName == "" {
		return ErrInvalidName
	} else {
		err := p.addPatient(&p.root, PatientName, doctorName, apptdt)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

/*----------------------------Search Operations--------------------------*/

//binary search in bst for node with matching name
func (p *Patients) SearchPatient(name string) (*Patient, error) {
	return p.searchPatNodes(p.root, name)
}

//search;
//binary search in bst for node with matching name
func (p *Patients) searchPatNodes(pat *Patient, name string) (*Patient, error) {
	if pat == nil {
		return &Patient{}, ErrPatientNotFound
	} else {
		if pat.Info.PatientName == name {
			return pat, nil
		} else {
			if name < pat.Info.PatientName {
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
func (p *Patients) printPatients(pat *Patient) string {
	if pat != nil {
		return p.printPatients(pat.left) + fmt.Sprintf(pat.Info.PatientName) + "<br>" + p.printPatients(pat.right)
	}
	return ""
}

//Prints out all patients' name
func (p *Patients) Print() string {
	return p.printPatients(p.root)
}
