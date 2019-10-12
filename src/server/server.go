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

func informationPageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/info.html"))

	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/info.html"))

	if r.Method != http.MethodPost {
		t := template.Must(template.ParseFiles("./templates/login.html"))
		t.Execute(w, nil)
		return
	}

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
	
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
