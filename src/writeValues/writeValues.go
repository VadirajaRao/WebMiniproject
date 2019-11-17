package writeValues

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// LogIn represents the credetials for database login.
type LogIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"dbname"`
}

// Signup represents the data for adding new user.
type Signup struct {
	Fname string
	Lname string
	Uname string
	Pwd string
	Rpwd string
	Mail string
}

// Product information representation
type Product struct {
	Pname string
	Ouid  int
	Luid  int
}

// extractCredentials extracts the login credentials from the JSON file.
func extractingCredentials() (LogIn, error) {
	// ReadAll is used to read login credentials from the JSON file.
	cred, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		return LogIn{}, errors.Wrap(err, "file read fail, passwords.json")
	}

	// credentials holds the login credentials.
	var credentials LogIn

	// Unmarshal is used to extract the credentials from the JSON file and save it
	// in credentials variable.
	err = json.Unmarshal(cred, &credentials)
	if err != nil {
		return LogIn{}, errors.Wrap(err, "failed to extract credentials from JSON")
	}

	return credentials, nil
}

// connectToDatabase uses the credentials to login to the database and returns a
// connection object to the database.
func connectingToDatabase(credentials LogIn) (*sql.DB, error) {
	// Open is used to open a connection to the database.
	db, err := sql.Open(
		"mysql",
		credentials.Username+":"+credentials.Password+"@(127.0.0.1:3306)/"+
			credentials.Database+"?parseTime=true",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create conenction")
	}

	return db, nil
}

func setup() (*sql.DB, error) {
	credentials, err := extractingCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract credentials")
	}

	db, err := connectingToDatabase(credentials)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return db, nil
}


// Function to make a new entry into the `USER` table which is used to store the
// information about a user.
func CreateUser(signupCred *Signup) error{
	db, err := setup()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user (fname, lname, name, mail, pwd) VALUES (?, ?, ?, ?, ?)
  `
	_, err = db.Exec(
		query,
		signupCred.Fname,
		signupCred.Lname,
		signupCred.Uname,
		signupCred.Mail,
		signupCred.Pwd,
	)

	if err != nil {
		return errors.Wrap(err, "failed to insert into database.")
	}

	return nil
}

// Function to make a new entry into the `PRODUCT` table which is used to store
// the information about a product
func CreateProduct(productCred *Product) error {
	db, err := setup()
	if err != nil {
		return err
	}

	query := "INSERT INTO product (pname, ouid, luid) VALUES (?, ?, ?)"

	_, err = db.Exec (
		query,
		productCred.Pname,
		productCred.Ouid,
		productCred.Luid,
	)

	if err != nil {
		return err
	}

	return nil
}

// Function to make a new entry into the `PRODUCT_BACKLOG` table which is used to
// store the information regarding the features requested by the owner.
func ProdLogEntry(pid int, issue string) (error) {
	db, err := setup()
	if err != nil {
		return err
	}

	query := "INSERT INTO product_backlog (pid, issue) VALUES (?, ?)"

	_, err = db.Exec(query, pid, issue)
	if err != nil {
		return err
	}

	return nil
}

// Function to remove an entry from the product backlog
func DroppingProdLog(pid int, issue string) (error) {
	db, err := setup()
	if err != nil {
		return err
	}

	query := "DELETE FROM product_backlog WHERE pid = ? AND issue = ?"

	_, err = db.Exec(query, pid, issue)
	if err != nil {
		return errors.Wrap(err, "failed to find the record")
	}

	return nil
}

// Function to add an entry into sprint backlog
func AddingToSprintLog(sid int, pid int, issue string, status string, uid int) (error) {
	db, err := setup()
	if err != nil {
		return err
	}

	query := `
    INSERT INTO sprint_backlog(sid, pid, issue, status, uid)
    VALUES (?, ?, ?, ?, ?)
  `

	_, err = db.Exec(query, sid, pid, issue, status, uid)
	if err != nil {
		return errors.Wrap(err, "unable to insert into sprint_cycle")
	}

	return nil
}

// Function to remove a log entry from the sprint_backlog table
func DroppingSprintLog(sid int, pid int, issue string) (error) {
	db, err := setup()
	if err != nil {
		return err
	}

	query := "DELETE FROM sprint_backlog WHERE sid = ? AND pid = ? AND issue = ?"

	_, err = db.Exec(query, sid, pid, issue)
	if err != nil {
		return errors.Wrap(err, "failed to find the error")
	}

	return nil
}
