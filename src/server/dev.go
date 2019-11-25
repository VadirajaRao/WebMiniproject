/*
 * This file contains all the actions performed when the active user is a
 * developer.
 */
package main

import (
	"log"
	"net/http"
	"html/template"

	"fetchValues"
	"writeValues"
)

type ErrMsg struct {
	Flag bool
	Msg  string
}

// Function to handle dev sprint backlog page
func devSprintBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_sprint_backlog.html"))

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the pid of the dev
	pid, err := fetchValues.FetchPIDDev(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the sid of the PID
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

// Function to handle dev progress handler
func devProgressHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_progress.html"))
	
	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the pid of the dev
	pid, err := fetchValues.FetchPIDDev(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the sid of the PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the in progress issues
	backlog, err := fetchValues.DevInProgressLog(sid, pid)
	if err != nil {
		if backlog.Msg == "No issue in-progress" {
			t.Execute(w, nil)
			return
		}

		log.Fatal(err)
	}

	t.Execute(w, backlog)
}

// Function to handle dev manage task handler
func devManageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_manage.html"))

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	uid := session.Values["user"].(int)

	// Fetching the pid of the dev
	pid, err := fetchValues.FetchPIDDev(uid) // Maybe the point of failure
	if err != nil {
		log.Fatal(err)
	}

	// Extracting Product name based on PID
	Pname, err := fetchValues.ExtractingProdName(pid)
	if err != nil {
		log.Fatal(err)
	}

	if r.Method != http.MethodPost {
		t.Execute(w, Pname)
		return
	}

	// Fetching the sid of the PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}

	// Reading form data
	issue := r.FormValue("feature")
	action := r.FormValue("action")

	// Updatin sprint_backlog
	if action == "progress" {
		err = writeValues.UpdatingInprogress(sid, pid, issue, "INPROGRESS", uid)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/dev/in-progress", http.StatusFound)
	} else {
		err = writeValues.UpdatingCompleted(sid, pid, issue, "COMPLETED", uid)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/dev/completed", http.StatusFound)
	}
}

// Function to handle dev completed handler
func devCompletedHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_completed.html"))

	// Extracting session information
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the pid of the dev
	pid, err := fetchValues.FetchPIDDev(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	// Fetching the sid of the PID
	sid, err := fetchValues.FetchingSID(pid)
	if err != nil {
		log.Fatal(err)
	}

	backlog, err := fetchValues.DevCompletedLog(sid, pid)
	if err != nil {
		if backlog.Msg == "No issue complete" {
			t.Execute(w, nil)
			return
		}

		log.Fatal(err)
	}

	t.Execute(w, backlog)
}
