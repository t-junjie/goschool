package bookingSystem

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type DoctorInfo struct {
	DoctorName   string
	Appointments Queue
	Availability []time.Time
}

//node
type Doctor struct {
	Info  DoctorInfo
	left  *Doctor
	right *Doctor
}

//bst
type Doctors struct {
	root *Doctor
}

var (
	ErrInvalidName = errors.New("Invalid Name")
	ErrInvalidList = errors.New("List is empty!")
	ErrDocNotFound = errors.New("The named doctor was not found!")
)

/*----------------------------Helper Methods--------------------------*/

func CreateDoctors() Doctors {
	return Doctors{nil}
}

/*----------------------------Insert Operations--------------------------*/

//iterate through list to add doctors to bst
func (d *Doctors) InitDoctors(list []string) error {
	if len(list) == 0 {
		return errors.Wrap(ErrInvalidList, "unable to initialize list of doctors ")
	}
	for _, name := range list {
		err := d.add(name)
		if err != nil {
			return err
		}
	}
	return nil
}

//insert;
//add individual node to bst
func (d *Doctors) addDoctor(doc **Doctor, name string) error {
	if *doc == nil {
		newDoctorInfo := DoctorInfo{name, Queue{nil, nil, 0}, initAvailability()}
		newNode := &Doctor{newDoctorInfo, nil, nil}
		*doc = newNode
		return nil
	}
	if name < (*doc).Info.DoctorName {
		err := d.addDoctor(&((*doc).left), name)
		if err != nil {
			return err
		}
	} else {
		err := d.addDoctor(&((*doc).right), name)
		if err != nil {
			return err
		}
	}
	return nil
}

//wrapper function for addDoctor
func (d *Doctors) add(name string) error {
	if name == "" {
		return errors.Wrap(ErrInvalidName, "name is provided as empty string")
	} else {
		err := d.addDoctor(&d.root, name)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

/*----------------------------Search Operations--------------------------*/

//binary search in bst for node with matching name
func (d *Doctors) SearchDoctor(name string) (*Doctor, error) {
	return d.searchDocNodes(d.root, name)
}

//search;
//binary search in bst for node with matching name
func (d *Doctors) searchDocNodes(doc *Doctor, name string) (*Doctor, error) {
	if doc == nil {
		return &Doctor{}, errors.Wrap(ErrDocNotFound, name)
	} else {
		if doc.Info.DoctorName == name {
			return doc, nil
		} else {
			if name < doc.Info.DoctorName {
				return d.searchDocNodes(doc.left, name)
			} else {
				return d.searchDocNodes(doc.right, name)
			}
		}
	}
}

/*----------------------------Print Operations--------------------------*/

//print;
//in order traversal for bst which prints out doctor's name
func (d *Doctors) printDoctors(doc *Doctor) string {
	if doc != nil {
		return d.printDoctors(doc.left) + fmt.Sprintf(doc.Info.DoctorName) + "<br>" + d.printDoctors(doc.right)
	}
	return ""
}

//wrapper for printDoctors
func (d *Doctors) Print() string {
	return d.printDoctors(d.root)
}
