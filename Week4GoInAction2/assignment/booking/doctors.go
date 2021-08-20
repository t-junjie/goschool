package booking

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

//DoctorInfo contains information regarding a doctor.
type DoctorInfo struct {
	DoctorName   string
	Appointments Queue
	Availability []time.Time
}

//Doctor is a Node in the Doctors BST.
type Doctor struct {
	Info  DoctorInfo
	left  *Doctor
	right *Doctor
}

//Doctors is a Binary Search Tree.
type Doctors struct {
	root *Doctor
}

/*----------------------------Helper Methods--------------------------*/

//CreateDoctors returns an empty Doctors BST.
func CreateDoctors() Doctors {
	return Doctors{nil}
}

/*----------------------------Insert Operations--------------------------*/

//InitDoctors iterates through a list of doctor names to add Doctor nodes to the Doctors BST.
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

//addDoctor creates a Doctor node and adds it to the Doctors BST.
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

//add is a convenience function that wraps addDoctor,
//requiring only that a new doctor's name is passed to it.
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

//SearchDoctor is a convenience function that wraps searchDocNodes,
//requiring only a doctor's name to be passed to it.
func (d *Doctors) SearchDoctor(name string) (*Doctor, error) {
	return d.searchDocNodes(d.root, name)
}

//searchDocNodes recursively searches the Doctors BST for a node with the specified name.
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

//printDoctors does an in-order traversal through the BST and returns the doctor's name.
func (d *Doctors) printDoctors(doc *Doctor) string {
	if doc != nil {
		return d.printDoctors(doc.left) + fmt.Sprintf(doc.Info.DoctorName) + "<br>" + d.printDoctors(doc.right)
	}
	return ""
}

//Print returns a formatted string that contains the names of all the doctors in the Doctors BST.
func (d *Doctors) Print() string {
	return d.printDoctors(d.root)
}
