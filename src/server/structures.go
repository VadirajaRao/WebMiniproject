package main

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

// This structure is used to pass message to the `login.html` page to indicate
// incorrect input.
type Msg struct {
	UFlag bool
	PFlag bool
	Flag bool
	Msg string
}

// This structure is used to pass message to the `signup.html` page to indicate
// incorrect input.
type SignMsg struct {
	MFlag bool
	Flag bool
}

// This structure is used to pass message to the `product.html` page to indicate
// incorrect input.
type ProductMsg struct {
	Oflag bool
	Lflag bool
}
