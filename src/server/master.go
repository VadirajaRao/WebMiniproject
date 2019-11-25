/*
 * This file contains all the actions user performed when the active user is
 * Scrum master.
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

var dummyUID int = 1

// Function to handle leader product backlog page
func masterProdBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_prod_backlog.html"))

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID based on the master UID
	pid, err := fetchValues.FetchPIDLeader(session.Values["user"].(int))
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

	t.Execute(w, backlog)
}

// Function to handle leader sprint backlog page
func masterSprintBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFiles("./templates/leader_sprint_backlog.html"),
	)

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID based on the master UID
	pid, err := fetchValues.FetchPIDLeader(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching active SID of the PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}
	
	// Fetching the sprint backlog
	backlog, err := fetchValues.FetchingSprintLog(sid, pid)
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, backlog)
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
	err = writeValues.AddingToSprintLog(sid, pid, issue, "AVAILABLE", dummyUID)
	if err != nil {
		log.Fatal(err) // Unable to add the entry into sprint backlog error.
	}

	// Redirecting to sprint backlog
	http.Redirect(w, r, "/master/sprint-backlog", http.StatusFound)
}

// Function to handle leader remove feature rpage
func masterRemoveFeatureHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_rem_feature.html"))

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID based on master UID
	pid, err := fetchValues.FetchPIDLeader(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching SID based on PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}

	// Reding the issue from the user
	issue := r.FormValue("feature")

	// Deleting the values from the sprint_backlog table
	err = writeValues.DroppingSprintLog(sid, pid, issue)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/master/sprint-backlog", http.StatusFound)
}

// Function to handle leader manage developer page
func masterManageDevHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader_manage.html"))

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching PID based on master UID
	pid, err := fetchValues.FetchPIDLeader(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Assigned developer list
	developerList, err := fetchValues.ExtractingDevs(pid)
	if err != nil {
		log.Fatal(err)
	}

	if r.Method != http.MethodPost {
		t.Execute(w, developerList)
		return
	}

	// Reading dev mail id input
	mail := r.FormValue("mail")

	// Extracting UID of the developer
	uid, err := fetchValues.FetchUID(mail)
	if err != nil {
		log.Fatal(err)
	}

	// Adding developer
	err = writeValues.AddingDev(pid, uid)
	if err != nil {
		log.Fatal(err)
	}
	
	t.Execute(w, developerList)
}
