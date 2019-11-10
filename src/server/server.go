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

// Session key initialization. (hide)
var store = sessions.NewCookieStore([]byte("MT-15vsR15"))

// Function to handle the main page.
// This is where the application begins
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
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

	t := template.Must(template.ParseFiles("./templates/index.html"))
	t.Execute(w, nil)
}

// Function to handle the information page
func informationPageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/info.html"))

	t.Execute(w, nil)
}

// Function to retrieve handle for redirection
func retrievingHandle(uid int) (string, error) {
	flag, err := auth.CheckIfOwner(uid)
	if err != nil {
		log.Fatal(err)
	}

	if flag == true {
		return "/owner", nil
	}

	flag, err = auth.CheckIfLeader(uid)
	if err != nil {
		log.Fatal(err)
	}

	if flag == true {
		return "/master", nil
	}

	return "/dev", nil
}

// Function to handle the login page
// This function verifies the user who wants to login to his/her account.
func loginHandler(w http.ResponseWriter, r *http.Request) {	
	session, err := store.Get(r, "session-name-1")
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.ParseFiles("./templates/login.html"))

	if r.Method != http.MethodPost {
		if session.Values["authenticated"] == true {
			http.Redirect(w, r, "/main", http.StatusFound)
		}
		t.Execute(w, nil)
		return
	}

	loginCred := Login {
		Usermail: r.FormValue("usermail"),
		Password: r.FormValue("password"),
	}

	val, err := fetchValues.LoginVerification(loginCred.Usermail)

	if err != nil && val != "NoUser" {
		log.Fatal(err)
	}

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

	session.Values["user"], err = fetchValues.FetchUID(loginCred.Usermail)
	if err != nil {
		log.Fatal(err)
	}
	session.Values["authenticated"] = true

	err = session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

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
	}

	err = writeValues.CreateUser(&signupCred)

	if err != nil {
		t.Execute(w, nil)
	}

	t = template.Must(template.ParseFiles("./templates/done.html"))
	t.Execute(w, nil)
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

// Function to handle owner page
func ownerHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner.html"))
	t.Execute(w, nil)
}

// Function to handle owner product backlog page
func ownerBacklogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_backlog.html"))
	t.Execute(w, nil)
}

// Function to handle owner adding feature page
func ownerAddHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_add.html"))
	t.Execute(w, nil)
}

// Function to handle owner remove feature page
func ownerRemHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/owner_remove.html"))
	t.Execute(w, nil)
}

// Function to hande leader page
func masterHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/leader.html"))
	t.Execute(w, nil)
}

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
	t.Execute(w, nil)
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

// Function to handle dev page
func devHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/dev.html"))
	t.Execute(w, nil)
}

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

	r.HandleFunc("/owner", ownerHandler)
	r.HandleFunc("/owner/backlog", ownerBacklogHandler)
	r.HandleFunc("/owner/add-feature", ownerAddHandler)
	r.HandleFunc("/owner/remove-feature", ownerRemHandler)

	r.HandleFunc("/master", masterHandler)
	r.HandleFunc("/master/prod-backlog", masterProdBacklogHandler)
	r.HandleFunc("/master/sprint-backlog", masterSprintBacklogHandler)
	r.HandleFunc("/master/add-feature", masterAddFeatureHandler)
	r.HandleFunc("/master/remove-feature", masterRemoveFeatureHandler)
	r.HandleFunc("/master/manage-developers", masterManageDevHandler)
	
	r.HandleFunc("/dev", devHandler)
	r.HandleFunc("/dev/sprint-backlog", devSprintBacklogHandler)
	r.HandleFunc("/dev/in-progress", devProgressHandler)
	r.HandleFunc("/dev/manage-task", devManageHandler)
	r.HandleFunc("/dev/completed", devCompletedHandler)

	r.HandleFunc("/logout", logoutHandler)
	
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
