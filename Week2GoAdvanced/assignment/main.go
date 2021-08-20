package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wgMain sync.WaitGroup

const initialPrompt string = `

Please choose whether to interact with the system as a user or admin:
1. Admin
2. User
`

func main() {

	runtime.GOMAXPROCS(2)
	var input int
	//initialize doctors
	docList := []string{"SNG, MATTHEW", "LUO, GLENN", "CHUA YI YOU", "YANG, CHELSEA", "FOO REN HAO"}
	doctors := Doctors{nil}
	patients := Patients{nil}
	/*
		We assume that all doctors' starting time is from workStartHour to workEndHour
		and that appointments are in 30 minutes slots, barring lunchHour.
		We further propose that appointment booking is only available from 7 days from now until 14 days from now.
		(as set in OneWeekFromNow/TwoWeeksFromNow variable in availability.go for easier debugging)
	*/
	err := doctors.initDoctors(docList)

	if err != nil {
		panic(err)
	}

	for {
		fmt.Println(initialPrompt)
		fmt.Scanln(&input)
		if input != 1 && input != 2 {
			fmt.Println(ErrInvalidOption.Error())
			fmt.Scanln(&input)
			time.Sleep(1000 * time.Millisecond)
		}
		wgMain.Add(1)
		go chooseInterface(input, &doctors, &patients)
		wgMain.Wait()
	}
}

func chooseInterface(input int, doc *Doctors, pat *Patients) {
	defer wgMain.Done()
	switch input {
	case 1:
		fmt.Println(adminPrompt)
		runAdminMenu(doc, pat)
	case 2:
		fmt.Println(userPrompt)
		runUserMenu(doc, pat)
	default:
		runUserMenu(doc, pat)
	}
	time.Sleep(1000 * time.Millisecond)
}
