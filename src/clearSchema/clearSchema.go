/*
Package clearSchema implements functions that is used to destroy the database
schema. This package has to be executed only to take down the schema and all the
data. Be careful while using this because all data will be lost.
*/
package clearSchema

import (
	"database/sql"
	"os"
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

func extractCredentials() (LogIn, error) {
	file, err := os.Open("./passwords.json")
	if err != nil {
		return LogIn{}, errors.Wrap(err, "file open fail, passwords.json")
	}

	cred, err := ioutil.ReadAll(file)
	if err != nil {
		return LogIn{}, errors.Wrap(err, "file read fail, passwords.json")
	}

	var credentials LogIn

	err = json.Unmarshal(cred, &credentials)
	if err != nil {
		return LogIn{}, errors.Wrap(err, "failed to extract credentials from JSON")
	}

	return credentials, nil
}

func connectToDatabase(credentials LogIn) (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		credentials.Username + ":" + credentials.Password + "@(127.0.0.1:3306)/" +
			credentials.Database + "?parseTime=true",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open a connection to the database")
	}

	return db, nil
}

func dropTables(db *sql.DB) error {
	tables := [...]string{
		"developer",
		"sprint_backlog",
		"sprint_cycle",
		"product_backlog",
		"product",
		"user",
	}

	var i int32

	for i < int32(len(tables)) {
		tableName := tables[i]
		query := "DROP TABLE IF EXISTS " + tableName + ";"

		_, err := db.Exec(query)
		if err != nil {
			return errors.Wrap(err, "failed to drop table " + tableName)
		}
	}

	return nil
}

func ClearTables() error {
	credentials, err := extractCredentials()
	if err != nil {
		return errors.Wrap(err, "failed to extract credentials")
	}

	db, err := connectToDatabase(credentials)
	if err != nil {
		return errors.Wrap(err, "failed to create new connection")
	}

	err = dropTables(db)
	if err != nil {
		return err
	}

	return nil
}
