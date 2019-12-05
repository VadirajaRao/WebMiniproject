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
	"auth"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Session key initialization.
var store = sessions.NewCookieStore([]byte("MT-15vsR15"))

// Function to handle the main page.
// This is where the application begins
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	// Fetching the session details for a particular session.
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Clearing the session by removing the user that was currently logged-in
	session.Values["user"] = -1 // This holds the user uid
	session.Values["authenticated"] = false // This tells if the session is active.

	// Saving the updated session information.
	err = session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

	// Rendering the template for the `index.html` page.
	t := template.Must(template.ParseFiles("./templates/index.html"))
	t.Execute(w, nil)
}

// Function to handle the information page
func informationPageHandler(w http.ResponseWriter, r *http.Request) {
	// Rendering the template for the `info.html` page.
	t := template.Must(template.ParseFiles("./templates/info.html"))
	t.Execute(w, nil)
}

/* Function to retrieve handle for redirection
 *
 * This function checks the user that is trying to login and redirects them to
 * appropriate landing page based on their role i.e., as a owner, scrum master or
 * a developer.
 *
 * The function will first check if the UID is of a product owner and then for
 * scrum master. If neither of them turn out to be true then it will redirect as
 * developer by default.
 */ 
func retrievingHandle(uid int) (string, error) {
	// `flag` is set if the UID of the user corresponds to a product owner.
	flag, err := auth.CheckIfOwner(uid)
	if err != nil {
		log.Fatal(err)
	}

	if flag == true {
		return "/owner/backlog", nil
	}

	// `flag` is ser if the UID of the user corresponds to a scrum master.
	flag, err = auth.CheckIfLeader(uid)
	if err != nil {
		log.Fatal(err)
	}

	if flag == true {
		return "/master/prod-backlog", nil
	}

	return "/dev/sprint-backlog", nil
}

/* Function to handle the login page
 *
 * This function verifies the user who wants to login to their account.
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieving session information for the user
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	// Preparing the tempplate for `login.html`
	t := template.Must(template.ParseFiles("./templates/login.html"))

	// Checking the HTTP method to verify if the form has been submitted or not.
	// This is necessary to render the form in case the session has not yet begun
	// or the user has not logged in yet.
	if r.Method != http.MethodPost {
		// To check if a session is already running w=in which case the user need not
		// login again.
		if session.Values["authenticated"] == true {
			handle, err := retrievingHandle(session.Values["user"].(int))
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, handle, http.StatusFound)
		}

		// Render the login form since session is inactive.
		t.Execute(w, nil)
		return
	}

	// Fetching the details submitted in the form by the user.
	loginCred := Login {
		Usermail: r.FormValue("usermail"),
		Password: r.FormValue("password"),
	}

	// Fetching the actual password associated with the mail id provided
	val, err := fetchValues.LoginVerification(loginCred.Usermail)
	if err != nil && val != "NoUser" {
		log.Fatal(err)
	}

	// Invalid mail-id.
	if val == "NoUser" {
		x := Msg {
			Flag: false,
			UFlag: true,
			PFlag: false,
			Msg: "Invalid email",
		}

		t.Execute(w, x)
		return
	}

	// Password mismatch.
	if val != loginCred.Password {
		x := Msg {
			UFlag: false,
			PFlag: true,
			Msg: "Incorrect Password",
			Flag: false,
		}

		t.Execute(w, x)
		return
	}

	// Starting the session after successful authorization
	session.Values["user"], err = fetchValues.FetchUID(loginCred.Usermail)
	if err != nil {
		log.Fatal(err)
	}
	session.Values["authenticated"] = true

	// Updating the session.
	err = session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieving the destination URL based on the role.
	handle, err := retrievingHandle(session.Values["user"].(int))
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, handle, http.StatusFound)
}

// Function to handle the sign up function. This function will create a new user
// by making entries into the database.
func signupHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/signup.html"))

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

	/*
  // This feature should be worked on. Something is not write.
	val, err := fetchValues.CheckMailInUser(signupCred.Mail)
	if err != nil {
		log.Fatal(err)
	}

	if val == true {
		x := SignMsg {
			MFlag: true,
			Flag: false,
		}
		
		t.Execute(w, x)
		return
	}*/

	err := writeValues.CreateUser(&signupCred)

	if err != nil {
		t.Execute(w, nil)
	}

	//t = template.Must(template.ParseFiles("./templates/done.html"))
	//t.Execute(w, nil)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Function to retrieve UID of the owner of the product
func retrievingOwnerUID(mail string, w http.ResponseWriter) (int, error) {
	t := template.Must(template.ParseFiles("./templates/new_product.html"))
	
	val, err := fetchValues.CheckOwnerMailInProduct(mail)
	if err != nil {
		log.Fatal(err)
	}
	if val == true {
		x := ProductMsg {
			Oflag: true,
			Lflag: false,
		}

		t.Execute(w,x)
		return 0, nil
	}

	uid, err := fetchValues.FetchUID(mail)

	return uid, err
}

// Function to retrieve UID of the leader of the product
func retrievingLeaderUID(mail string, w http.ResponseWriter) (int, error) {
	t := template.Must(template.ParseFiles("./templates/new_product.html"))
	
	val, err := fetchValues.CheckLeaderMailInProduct(mail)
	if err != nil {
		log.Fatal(err)
	}

	if val == true {
		x := ProductMsg {
			Oflag: false,
			Lflag: true,
		}

		t.Execute(w, x)
		return 0, nil
	}

	uid, err := fetchValues.FetchUID(mail)

	return uid, err
}

// Function to handle `new_product.html` or creation of new product page.
func newProductHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/new_product.html"))
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}
	productCred := ProductDetails {
		Pname: r.FormValue("pname"),
		Omail: r.FormValue("omail"),
		Lmail: r.FormValue("lmail"),
	}
	var finProdCred writeValues.Product
	finProdCred.Pname = productCred.Pname

	var err error

	finProdCred.Ouid, err = retrievingOwnerUID(productCred.Omail, w)
	if err != nil {
		log.Fatal(err)
	}

	finProdCred.Luid, err = retrievingLeaderUID(productCred.Lmail, w)
	if err != nil {
		log.Fatal(err)
	}

	err = writeValues.CreateProduct(&finProdCred)
	if err != nil {
		t.Execute(w, nil)
	}

	t = template.Must(template.ParseFiles("./templates/done.html"))
	t.Execute(w, nil)
}

// Function to logout an user
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	session.Values["user"] = -1
	session.Values["authenticated"] = false

	err = session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
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

	//r.HandleFunc("/owner", ownerHandler)
	r.HandleFunc("/owner/backlog", ownerBacklogHandler)
	r.HandleFunc("/owner/add-feature", ownerAddHandler)
	r.HandleFunc("/owner/remove-feature", ownerRemHandler)

	//r.HandleFunc("/master", masterHandler)
	r.HandleFunc("/master/prod-backlog", masterProdBacklogHandler)
	r.HandleFunc("/master/sprint-backlog", masterSprintBacklogHandler)
	r.HandleFunc("/master/add-feature", masterAddFeatureHandler)
	r.HandleFunc("/master/remove-feature", masterRemoveFeatureHandler)
	r.HandleFunc("/master/manage-developers", masterManageDevHandler)
	
	//r.HandleFunc("/dev", devHandler)
	r.HandleFunc("/dev/sprint-backlog", devSprintBacklogHandler)
	r.HandleFunc("/dev/in-progress", devProgressHandler)
	r.HandleFunc("/dev/manage-task", devManageHandler)
	r.HandleFunc("/dev/completed", devCompletedHandler)

	r.HandleFunc("/logout", logoutHandler)
	
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
