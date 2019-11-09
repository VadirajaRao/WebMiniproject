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
	Usermail string
	Password string
}

// This structure is used to store the sign up information submitted by the user
// in the signup.html form.
type Signup struct {
	Fname string
	Lname string
	Uname string
	Pwd string
	Rpwd string
	Mail string
}

// This structure is used to store the new product information submitted by the
// user in the new_product.html form.
type ProductDetails struct {
	Pname string
	Omail string
	Lmail string
}

type Msg struct {
	UFlag bool
	PFlag bool
	Flag bool
	Msg string
}

type SignMsg struct {
	MFlag bool
	Flag bool
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

func newProductHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./templates/new_product.html"))
	//t.Execute(w, nil)

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	type ProductMsg struct {
		Oflag bool
		Lflag bool
	}

	productCred := ProductDetails {
		Pname: r.FormValue("pname"),
		Omail: r.FormValue("omail"),
		Lmail: r.FormValue("lmail"),
	}

	var finProdCred writeValues.Product

	finProdCred.Pname = productCred.Pname

	val, err := fetchValues.CheckOwnerMailInProduct(productCred.Omail)
	if err != nil {
		log.Fatal(err)
	}

	if val == true {
		x := ProductMsg {
			Oflag: true,
			Lflag: false,
		}

		t.Execute(w, x)
		return
	}

	finProdCred.Ouid, err = fetchValues.FetchUID(productCred.Omail)
	if err != nil {
		log.Fatal(err)
	}

	val, err = fetchValues.CheckLeaderMailInProduct(productCred.Lmail)
	if err != nil {
		log.Fatal(err)
	}

	if val == true {
		x := ProductMsg {
			Oflag: false,
			Lflag: true,
		}

		t.Execute(w, x)
		return
	}

	finProdCred.Luid, err = fetchValues.FetchUID(productCred.Lmail)
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
