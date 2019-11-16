/*
 * This file contains all the actions user performed when the active user is
 * Scrum master.
 */
package main

import (
	"log"
	"net/http"
	"html/template"

	"writeValues"
	"fetchValues"
)

// Function to handle leader product backlog page
func masterProdBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_prod_backlog.html"))
	t.Execute(w, nil)
}

// Function to handle leader sprint backlog page
func masterSprintBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFiles("./templates/leader_sprint_backlog.html"),
	)
	t.Execute(w, nil)
}

// Function to handle leader add feature page
func masterAddFeatureHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_add_feature.html"))
	//t.Execute(w, nil)

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID of the owner
	pid, err := fetchValues.FetchPIDLeader(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching active SID of the PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}

	// Reading the feature sent by the user
	issue := r.FormValue("feature")

	// Inserting the issue into the sprint backlog
	// The UID that is used here is a dummy user. Should fix the schema.
	err = writeValues.AddingToSprintLog(sid, pid, issue, "AVAILABLE", 20)
	if err != nil {
		log.Fatal(err) // Unable to add the entry into sprint backlog error.
	}

	// Redirecting to sprint backlog
	http.Redirect(w, r, "/master/sprint-backlog", http.StatusFound)
}

// Function to handle leader remove feature rpage
func masterRemoveFeatureHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_rem_feature.html"))
	t.Execute(w, nil)
}

// Function to handle leader manage developer page
func masterManageDevHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_manage.html"))
	t.Execute(w, nil)
}
