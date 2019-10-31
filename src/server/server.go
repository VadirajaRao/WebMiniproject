package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"html/template"

	"createSchema"
	"clearSchema"
	"writeValues"
	"fetchValues"

	"github.com/gorilla/mux"
)

// This structure is used to store the login information submitted by the user
// in the login.html form.
type Login struct {
	Username string
	Password string
}

type Signup struct {
	Fname string
	Lname string
	Uname string
	Pwd string
	Rpwd string
	Mail string
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

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	loginCred := Login {
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	val, err := fetchValues.LoginVerification(loginCred.Username)
	if err != nil && val != "NoUser" {
		log.Fatal(err)
	}

	if val == "NoUser" {
		type UserMsg struct {
			UFlag bool
			Flag bool
			PFlag bool
			UMsg string
		}

		x := UserMsg {
			Flag: false,
			UFlag: true,
			PFlag: false,
			UMsg: "Invalid Username",
		}

		t.Execute(w, x)
		return
	}

	if val != loginCred.Password {
		type PwdMsg struct {
			UFlag bool
			PFlag bool
			PMsg string
			Flag bool
		}
		
		x := PwdMsg {
			UFlag: false,
			PFlag: true,
			PMsg: "Incorrect Password",
			Flag: false,
		}

		t.Execute(w, x)
		return
	}

	t = template.Must(template.ParseFiles("./templates/done.html"))
	t.Execute(w, nil)
}

// Function to handle the sign up function. This function will create a new user
// by making entries into the database.
func signupHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/signup.html"))
	//t.Execute(w, nil)

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	signupCred := writeValues.Signup {
		Fname: r.FormValue("fname"),
		Lname: r.FormValue("lname"),
		Uname: r.FormValue("uname"),
		Pwd: r.FormValue("pwd"),
		Rpwd: r.FormValue("rpwd"),
		Mail: r.FormValue("mail"),
	}

	err := writeValues.CreateUser(&signupCred)

	if err != nil {
		t.Execute(w, nil)
	}

	t = template.Must(template.ParseFiles("./templates/done.html"))
	t.Execute(w, nil)
}

func newProductHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/new_product.html"))
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
	r.HandleFunc("/newProduct", newProductHandler)
	
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
