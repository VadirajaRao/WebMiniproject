package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"html/template"

	"createSchema"
	"clearSchema"

	"github.com/gorilla/mux"
)

// This structure is used to store the login information submitted by the user
// in the login.html form.
type Login struct {
	Username string
	Password string
}

// Function to handle the main page.
// This is where the application begins
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/index.html"))

	t.Execute(w, nil)
}

// Function to handle the information page
func informationPageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/info.html"))

	t.Execute(w, nil)
}

// Function to handle the login page
// This function verifies the user who wants to login to his/her account.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/login.html"))

	// If the HTTP request is using GET method, then the login form will be loaded
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}
	// The following section will be executed once the form has been submitted.
	// Since the form uses POST method to submit the form content, the above
	// section will not be executed.

	// This statement reads the value once the form is submitted.
	name := r.FormValue("username")

	// The following structure is used to send data back to the html page to show
	// a particular message. These values are made use in the if-else block in
	// login.html
	type Variables struct {
		Success bool
		Name string
	}

	x := Variables {
		Success: true,
			Name: name,
		}

	t.Execute(w, x)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/signup.html"))
	t.Execute(w, nil)
}

func main() {
	// Checking for arguments passed from the command line
	if len(os.Args) == 2 {
		// To create the database schema
		if os.Args[1] == "create" {
			err := createSchema.CreatingTables()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("It started working again")
			}
			return
		} else if os.Args[1] == "clear" {
			// To clear the database schema
			err := clearSchema.ClearingTables()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Clearing database complete")
			}
			return
		}
	}
	
	r := mux.NewRouter()

	r.HandleFunc("/", indexPageHandler)
	r.HandleFunc("/info", informationPageHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/signup", signupHandler)
	
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
