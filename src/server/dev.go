/*
 * This file contains all the actions performed when the active user is a
 * developer.
 */
package main

import (
	"net/http"
	"html/template"
)

/*
// Function to handle dev page
func devHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev.html"))
	t.Execute(w, nil)
}
*/

// Function to handle dev sprint backlog page
func devSprintBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_sprint_backlog.html"))
	t.Execute(w, nil)
}

// Function to handle dev progress handler
func devProgressHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_progress.html"))
	t.Execute(w, nil)
}

// Function to handle dev manage task handler
func devManageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_manage.html"))
	t.Execute(w, nil)
}

// Function to handle dev completed handler
func devCompletedHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev_completed.html"))
	t.Execute(w, nil)
}
