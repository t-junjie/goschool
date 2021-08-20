package routing

import (
	"html/template"
	"net/http"

	system "Assignment3/bookingSystem"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type userInfo struct {
	Username  string
	Password  []byte
	FirstName string
	LastName  string
}

type Result struct {
	User                userInfo
	DoctorsAvailability template.HTML
	DoctorAvailability  template.HTML
	BookingStatus       bool
	PatientsName        template.HTML
	PatientHistory      template.HTML
	DoctorsName         template.HTML
	DoctorAppointments  template.HTML
	DoctorName          template.HTML
	PatientName         template.HTML
}

var (
	doctors     system.Doctors
	patients    system.Patients
	mapUsers    = map[string]userInfo{}
	mapSessions = map[string]string{}
	tpl         *template.Template
	r           Result
)

func init() {
	doctors = system.CreateDoctors()
	patients = system.CreatePatients()
	docList := []string{"ALAN CHEUNG", "ANG BENG CHONG", "CHUA YI YOU", "KARLSSON BENGT GUNNAR", "FOO REN HAO"}
	//InitDoctors only returns an error for empty slice passed to it, so we can safely ignore the error
	_ = doctors.InitDoctors(docList)

	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mapUsers["admin"] = userInfo{"admin", bPassword, "admin", "admin"}

	tpl = template.Must(template.ParseGlob("templates/*"))
}

func Route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/home", user)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/booking", bookNewAppointment)
	http.HandleFunc("/searchAvailableDoctors", searchAvailableDoctors)
	http.HandleFunc("/searchDoctorAvailability", searchDoctorAvailability)
	http.HandleFunc("/viewDoctorAppointments", viewDoctorAppointments)
	http.HandleFunc("/choosePatient", choosePatient)
	http.HandleFunc("/chooseDoctor", chooseDoctor)
	http.HandleFunc("/viewPatientAppointmentHistory", viewPatientAppointmentHistory)
	http.HandleFunc("/searchAppointment", searchAppointment)
	http.HandleFunc("/editAppointment", editAppointment)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":5221", nil)
}

//Helper Functions

func isLoggedIn(req *http.Request) bool {
	loginCookie, err := req.Cookie("loginCookie")
	if err != nil {
		return false
	}
	username := mapSessions[loginCookie.Value]
	_, ok := mapUsers[username]
	return ok
}

func getUser(res http.ResponseWriter, req *http.Request) userInfo {
	// get current session cookie
	loginCookie, err := req.Cookie("loginCookie")
	if err != nil {
		id := uuid.NewV4()
		loginCookie = &http.Cookie{
			Name:  "loginCookie",
			Value: id.String(),
		}

	}
	http.SetCookie(res, loginCookie)

	// if the user exists already, get user
	var myUser userInfo
	if username, ok := mapSessions[loginCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

func isAdmin(res http.ResponseWriter, req *http.Request) bool {
	myUser := getUser(res, req)
	return myUser.Username == "admin"
}
