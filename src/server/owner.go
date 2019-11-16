/*
 * This file contains all the functions performed when the active user is a
 * product owner.
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"writeValues"
	"fetchValues"
)


// Function to handle owner product backlog page
func ownerBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_backlog.html"))
	/*
  t.Execute(w, nil)

	if r.Method != http.MethoPost {
		t.Execute(w, nil)
		return
	}
  */

	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID of the product
	pid, err := fetchValues.FetchPID(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the product backlog
	// Take care when the product has no logs initially.
	// Maybe redirect to add feature page directly in case there is no entry.
	backlog, err := fetchValues.FetchingProdLog(pid)
	if err != nil {
		log.Fatal(err)
	}

	for _, log := range backlog {
		fmt.Println(log.Feature)
	}

	t.Execute(w, backlog)
}

// Function to handle owner adding feature page
func ownerAddHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_add.html"))

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the PID of the product.
	pid, err := fetchValues.FetchPID(session.Values["user"].(int))
	// Unable to find the product corresponding to the user. Handle the error.
	if err != nil {
		log.Fatal(err)
	}

	// Taking the feature submitted by the user.
	// In case the priority level has to be added create a structure for this.
	issue := r.FormValue("feature")

	// Adding the feature into the database.
	// Create an error handler saying unable to add feature.
	err = writeValues.ProdLogEntry(pid, issue)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/owner/backlog", http.StatusFound)
}

// Function to handle owner remove feature page
func ownerRemHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_remove.html"))

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID of the owner
	pid, err := fetchValues.FetchPID(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Taking the feature to be deleted
	issue := r.FormValue("feature")

	// Deleting the feature from the backlog
	// Create a modification package
	err = writeValues.DroppingProdLog(pid, issue)
	if err != nil {
		log.Fatal(err) // unable to find the record
	}

	http.Redirect(w, r, "/owner/backlog", http.StatusFound)
}
