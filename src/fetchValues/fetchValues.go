package fetchValues

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type LogIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"dbname"`
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

// Function to create connection to the database before any operation is
// performed.
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

// Function to check if the mail id is valid and then returns the password
// associated with that mail id.
func LoginVerification(usermail string) (string, error) {
	db, err := setup()
	if err != nil {
		return "", err
	}

	var (
		uid int
		fname string
		lname string
		name string
		mail string
		pwd string
	)
	query := `SELECT * FROM user WHERE mail = ?`

	err = db.QueryRow(query, usermail).Scan(
		&uid, &fname, &lname, &name, &mail, &pwd,
	)
	if err != nil {
		return "NoUser", errors.Wrap(err, "query execution failed")
	}

	return pwd, nil
}

// Function to check if the mail id already exists in the database and returns a
// boolean value. The presence of user is checked in `user` table.
func CheckMailInUser(usermail string) (bool, error) {
	db, err := setup()
	if err != nil {
		return false, err
	}

	var (
		uid int
		fname string
		lname string
		name string
		mail string
		pwd string
	)

	query := "SELECT * FROM user WHERE mail = ?"

	err = db.QueryRow(query, usermail).Scan(
		&uid, &fname, &lname, &name, &mail, &pwd,
	)
	if err != nil {
		return false, errors.Wrap(err, "query execution failed")
	}

	return true, nil
}

// Function to return the UID of a particular mail, if the mail id exists in the
// database.
func FetchUID(usermail string) (int, error) {
	db, err := setup()
	if err != nil {
		return 0, err
	}

	var uid int

	query := "SELECT uid FROM user WHERE mail = ?"

	err = db.QueryRow(query, usermail).Scan(&uid)
	if err != nil {
		return 0, errors.Wrap(err, "unable to retrieve UID")
	}

	return uid, nil
}

// Function to check if the mail id already exists in the database and returns a
// boolean value. The presence of user is checked in `product` table.
func CheckOwnerMailInProduct(usermail string) (bool, error) {
	db, err := setup()
	if err != nil {
		return false, err
	}

	uid, err := FetchUID(usermail)
	if err != nil {
		return false, err
	}

	var pid int

	query := "SELECT pid FROM product WHERE ouid = ?"
	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return false, errors.Wrap(err, "query execution failed")
	}

	return true, nil
}

// Function to check if the mail id already exists in the database and returns a
// boolean value. The presence of user is checked in `product` table.
func CheckLeaderMailInProduct(usermail string) (bool, error) {
	db, err := setup()
	if err != nil {
		return false, err
	}

	uid, err := FetchUID(usermail)
	if err != nil {
		return false, err
	}

	var pid int

	query := "SELECT pid FROM product WHERE luid = ?"
	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return false, errors.Wrap(err, "query execution failed")
	}

	return true, nil
}
