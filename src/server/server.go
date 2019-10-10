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

// Function to handle the main page.
// This is where the application begins
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal("failed to parse the file")
	}

	t.Execute(w, nil)
}

// Function to handle information page
func informationPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/info.html")
	if err != nil {
		log.Fatal("failed to parse info")
	}

	t.Execute(w, nil)
}

// Function to handle login page.
func afterLogin(w http.ResponseWriter, r *http.Request) {
	const u string = "hello"
	const p string = "1234"
	
	t, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		log.Fatal("failed to parse login")
	}
	t.Execute(w, nil)

	/*
	r.ParseForm()

	if u != r.FormValue("username") || p != r.FormValue("password") {
		log.Fatal("Total breakdown")
	}

	t.Execute(w, nil)
  */
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

	// fs := http.FileServer(http.Dir("static/css"))
	// http.Handle("/static/css", http.StripPrefix("/static/css", fs))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	r := mux.NewRouter()

	r.HandleFunc("/", indexPageHandler)
	r.HandleFunc("/info", informationPageHandler)
	r.HandleFunc("/login", afterLogin)

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
