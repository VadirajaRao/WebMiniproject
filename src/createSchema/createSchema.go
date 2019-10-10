/*
Package createSchema implements a set of functions that are used to create the
schema required for the application. This particular pacakge has to be executed
only once, which is to setup the application on the system, because once executed
the database schema will already exist in the database.

If executed second time on the same system, there will be an error because of the
existence of the schema.
*/
package createSchema

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

// CreateTables creates the tables for the application after extracting the login
// credentials and establishing connection to the database.
func CreatingTables() error {
	// extractCredentials returns login credentials.
	credentials, err := extractingCredentials()
	if err != nil {
		return errors.Wrap(err, "failed to extract credentials")
	}

	// connectToDatabase returns a connection to the database.
	db, err := connectingToDatabase(credentials)
	if err != nil {
		return errors.Wrap(err, "failed to connect to the database")
	}

	// USER Table
	query := `
    CREATE TABLE user (
      uid INT AUTO_INCREMENT,
      fname VARCHAR(100) NOT NULL,
      lname VARCHAR(100) NOT NULL,
      name VARCHAR(100) NOT NULL,
      pwd VARCHAR(100) NOT NULL,
      PRIMARY KEY(uid)
    );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create USER table")
	}

	// PRODUCT Table
	query = `
    CREATE TABLE product (
      pid INT AUTO_INCREMENT,
      pname VARCHAR(100) NOT NULL,
      ouid INT NOT NULL,
      luid INT NOT NULL,
      PRIMARY KEY(pid),
      CONSTRAINT uid1_fk FOREIGN KEY(ouid) REFERENCES user(uid)
        ON DELETE CASCADE,
      CONSTRAINT uid2_fk FOREIGN KEY(luid) REFERENCES user(uid) ON DELETE CASCADE
  );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create PRODUCT table")
	}

	// DEVELOPER Table
	query = `
    CREATE TABLE developer (
      pid INT NOT NULL,
      uid INT NOT NULL,
      PRIMARY KEY(pid, uid),
      CONSTRAINT uid3_fk FOREIGN KEY(uid) REFERENCES user(uid) ON DELETE CASCADE,
      CONSTRAINT pid1_fk FOREIGN KEY(pid) REFERENCES product(pid)
        ON DELETE CASCADE
    );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create DEVELOPER table")
	}

	// PRODUCT_BACKLOG Table
	query = `
    CREATE TABLE product_backlog (
      pid INT NOT NULL,
      issue VARCHAR(500) NOT NULL,
      PRIMARY KEY (pid, issue),
      CONSTRAINT pid2_fk FOREIGN KEY(pid) REFERENCES product(pid)
        ON DELETE CASCADE
    );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create PRODUCT_BACKLOG table")
	}

	// SPRINT_CYCLE Table
	query = `
    CREATE TABLE sprint_cycle (
      sid INT AUTO_INCREMENT,
      pid INT NOT NULL,
      sdate DATE NOT NULL,
      edate DATE NOT NULL,
      PRIMARY KEY(sid, pid),
      CONSTRAINT pid3_fk FOREIGN KEY(pid) REFERENCES product(pid)
        ON DELETE CASCADE
    );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create SPRINT_CYCLE table")
	}

	// SPRINT_BACKLOG Table
	query = `
    CREATE TABLE sprint_backlog (
      sid INT NOT NULL,
      pid INT NOT NULL,
      issue VARCHAR(500) NOT NULL,
      status VARCHAR(25) NOT NULL DEFAULT 'unassigned',
      uid INT NOT NULL,
      PRIMARY KEY(sid, pid, issue),
      CONSTRAINT sid1_fk FOREIGN KEY(sid) REFERENCES sprint_cycle(sid)
        ON DELETE CASCADE,
      CONSTRAINT issue1_fk FOREIGN KEY(pid, issue) REFERENCES
        product_backlog(pid, issue) ON DELETE CASCADE,
      CONSTRAINT uid4_fk FOREIGN KEY(uid) REFERENCES user(uid)
        ON DELETE CASCADE
    );`

	_, err = db.Exec(query)

	if err != nil {
		return errors.Wrap(err, "failed to create SPRINT_BACKLOG table")
	}

	return nil
}
