package booking

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

//PatientInfo contains information regarding a patient.
type PatientInfo struct {
	PatientName        string
	AppointmentHistory Queue
}

//Patient is a Node in the Patients BST.
type Patient struct {
	Info  PatientInfo
	left  *Patient
	right *Patient
}

//Patients is a Binary Search Tree.
type Patients struct {
	root *Patient
}

/*----------------------------Helper Functions--------------------------*/

//CreatePatients returns an empty Patients BST.
func CreatePatients() Patients {
	return Patients{nil}
}

/*----------------------------Insert Operations--------------------------*/

//addPatient creates a Patient node and adds it to the Patients BST.
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

//Add inserts a Patient node to the Patients BST.
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

//SearchPatient is a convenience function that wraps searchPatNodes,
//requiring only a doctor's name to be passed to it.
func (p *Patients) SearchPatient(name string) (*Patient, error) {
	return p.searchPatNodes(p.root, name)
}

//searchPatNodes recursively searches the Patients BST for a node with the specified name.
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

//printPatients does an in-order traversal through the BST and returns the patients's name.
func (p *Patients) printPatients(pat *Patient) string {
	if pat != nil {
		return p.printPatients(pat.left) + fmt.Sprintf(pat.Info.PatientName) + "<br>" + p.printPatients(pat.right)
	}
	return ""
}

//Print returns a formatted string that contains the names of all the patients in the Patients BST.
func (p *Patients) Print() string {
	return p.printPatients(p.root)
}
