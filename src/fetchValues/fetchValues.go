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

type Backlog struct {
	ProdName string
	Feature []string
}

type ProductLog struct {
	Feature []string
	UID     []int
}

type SingleLog struct {
	Feature string
	Uname   string
}

type ProgressLog struct {
	Msg      string
	ProdName string
	Logs     []SingleLog
}

type DevList struct {
	Msg      string
	ProdName string
	Dev      []string
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

// Function to extract product name based on pid
func ExtractingProdName(pid int) (string, error) {
	db, err := setup()
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := "SELECT pname FROM product WHERE pid = ?"

	var pname string

	err = db.QueryRow(query, pid).Scan(&pname)
	if err != nil {
		return "", errors.Wrap(err, "unable to extract product name")
	}

	return pname, nil
}

// Function to check if the mail id is valid and then returns the password
// associated with that mail id.
func LoginVerification(usermail string) (string, error) {
	db, err := setup()
	if err != nil {
		return "", err
	}
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

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

// Function to retrieve the PID based on the owner UID
func FetchPID(uid int) (int, error) {
	var pid int
	
	db, err := setup()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	query := "SELECT pid FROM product WHERE ouid = ?"

	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return -1, errors.Wrap(err, "failing to fetch PID")
	}

	return pid, nil
}

// Function to retrieve the product backlog based on PID
func FetchingProdLog(pid int) (Backlog, error) {
	db, err := setup()
	if err != nil {
		return Backlog{}, err
	}
	defer db.Close()

	query := "SELECT issue FROM product_backlog WHERE pid = ?"

	rows, err := db.Query(query, pid)
	if err != nil {
		return Backlog{}, err
	}
	defer rows.Close()

	var backlog Backlog

	for rows.Next() {
		var log string

		err := rows.Scan(&log)
		if err != nil {
			return Backlog{}, errors.Wrap(err, "failed to process a row")
		}

		backlog.Feature = append(backlog.Feature, log)
	}

	err = rows.Err()
	if err != nil {
		return Backlog{}, errors.Wrap(err, "failed after processing rows")
	}

	backlog.ProdName, err = ExtractingProdName(pid)
	if err != nil {
		return Backlog{}, err
	}

	return backlog, nil
}

// Function to retrieve the PID based on the master UID
func FetchPIDLeader(uid int) (int, error) {
	var pid int
	
	db, err := setup()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	query := "SELECT pid FROM product WHERE luid = ?"

	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return -1, errors.Wrap(err, "failing to fetch PID")
	}

	return pid, nil
}

// Fetching SID given PID from the SPRINT_CYCLE table
func FetchingSID(pid int) (int, error) {
	var sid int
	
	db, err := setup()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	query := "SELECT sid FROM sprint_cycle WHERE pid = ?"

	err = db.QueryRow(query, pid).Scan(&sid)
	if err != nil {
		return -1, errors.Wrap(err, "failing to fetch sid")
	}

	return sid, nil
}

// Function to fetch sprint backlog based on PID and SID
func FetchingSprintLog(sid int, pid int) (Backlog, error) {
	db, err := setup()
	if err != nil {
		return Backlog{}, err
	}
	defer db.Close()

	query := "SELECT issue FROM sprint_backlog WHERE sid = ? AND pid = ?"

	rows, err := db.Query(query, sid, pid)
	if err != nil {
		return Backlog{}, err
	}
	defer rows.Close()

	var backlog Backlog

	for rows.Next() {
		var log string

		err := rows.Scan(&log)
		if err != nil {
			return Backlog{}, errors.Wrap(err, "failed to process a row")
		}

		backlog.Feature = append(backlog.Feature, log)
	}

	err = rows.Err()
	if err != nil {
		return Backlog{}, errors.Wrap(err, "failed after processing rows")
	}

	// Extracting product name based on PID
	backlog.ProdName, err = ExtractingProdName(pid)
	if err != nil {
		return Backlog{}, err
	}

	return backlog, nil
}

// Function to fetch PID for Dev
func FetchPIDDev(uid int) (int, error) {
	var pid int
	
	db, err := setup()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	query := "SELECT pid FROM developer WHERE uid = ?"

	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return -1, errors.Wrap(err, "unable to find PID")
	}

	return pid, nil
}

// Function to extract username based on UID
func extractingUsername(uid int) (string, error) {
	db, err := setup()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var uname string
	query := "SELECT name FROM user WHERE uid = ?"

	err = db.QueryRow(query, uid).Scan(&uname)
	if err != nil {
		return "", errors.Wrap(err, "unable to extract user name")
	}

	return uname, nil
}

// Function to fetch the features in progress
func DevInProgressLog(sid int, pid int) (ProgressLog, error) {
	db, err := setup()
	if err != nil {
		return ProgressLog{}, err
	}
	defer db.Close()

	query := `
		SELECT issue, uid
    FROM sprint_backlog
    WHERE sid = ? AND pid = ? AND status="INPROGRESS"
  `

	rows, err := db.Query(query, sid, pid)
	if err != nil {
		return ProgressLog{}, err
	}
	defer rows.Close()

	var progressLog ProgressLog

	// Extracting product name based on PID
	progressLog.ProdName, err = ExtractingProdName(pid)
	if err != nil {
		return ProgressLog{}, err
	}

	if rows.Next() {
		var pLog SingleLog
		var uid int

		err := rows.Scan(&pLog.Feature, &uid)
		if err != nil {
			return ProgressLog{}, errors.Wrap(err, "failed to process a row")
		}

		pLog.Uname, err = extractingUsername(uid)
		if err != nil {
			return ProgressLog{}, err
		}

		progressLog.Logs = append(progressLog.Logs, pLog)
	} else {
		var pLog ProgressLog

		// Extracting product name based on PID
		pLog.ProdName, err = ExtractingProdName(pid)
		if err != nil {
			return ProgressLog{}, err
		}

		pLog.Msg = "No issue in-progress"

		return pLog, err
	}

	for rows.Next() {
		var pLog SingleLog
		var uid int

		err := rows.Scan(&pLog.Feature, &uid)
		if err != nil {
			return ProgressLog{}, errors.Wrap(err, "failed to process a row")
		}

		pLog.Uname, err = extractingUsername(uid)
		if err != nil {
			return ProgressLog{}, err
		}

		progressLog.Logs = append(progressLog.Logs, pLog)
	}
	
	err = rows.Err()
	if err != nil {
		return ProgressLog{}, errors.Wrap(err, "failed after processing rows")
	}

	progressLog.Msg = ""
	return progressLog, nil
}

// Function to fetch values from sprint-backlog table with COMPLETE status
func DevCompletedLog(sid int, pid int) (ProgressLog, error) {
	db, err := setup()
	if err != nil {
		return ProgressLog{}, err
	}
	defer db.Close()

	query := `
		SELECT issue, uid
    FROM sprint_backlog
    WHERE sid = ? AND pid = ? AND status="COMPLETED"
  `

	rows, err := db.Query(query, sid, pid)
	if err != nil {
		return ProgressLog{}, err
	}
	defer rows.Close()

	var progressLog ProgressLog

	// Extracting product name based on PID
	progressLog.ProdName, err = ExtractingProdName(pid)
	if err != nil {
		return ProgressLog{}, err
	}

	if rows.Next() {
		var pLog SingleLog
		var uid int

		err := rows.Scan(&pLog.Feature, &uid)
		if err != nil {
			return ProgressLog{}, errors.Wrap(err, "failed to process a row")
		}

		pLog.Uname, err = extractingUsername(uid)
		if err != nil {
			return ProgressLog{}, err
		}

		progressLog.Logs = append(progressLog.Logs, pLog)
	} else {
		var pLog ProgressLog

		// Extracting product name based on PID
		pLog.ProdName, err = ExtractingProdName(pid)
		if err != nil {
			return ProgressLog{}, err
		}

		pLog.Msg = "No issue complete"

		return pLog, err
	}

	for rows.Next() {
		var pLog SingleLog
		var uid int

		err := rows.Scan(&pLog.Feature, &uid)
		if err != nil {
			return ProgressLog{}, errors.Wrap(err, "failed to process a row")
		}

		pLog.Uname, err = extractingUsername(uid)
		if err != nil {
			return ProgressLog{}, err
		}

		progressLog.Logs = append(progressLog.Logs, pLog)
	}
	
	err = rows.Err()
	if err != nil {
		return ProgressLog{}, errors.Wrap(err, "failed after processing rows")
	}

	progressLog.Msg = ""
	return progressLog, nil
}

// Function to extract developers
func ExtractingDevs(pid int) (DevList, error) {
	db, err := setup()
	if err != nil {
		return DevList{}, err
	}
	defer db.Close()

	query := "SELECT uid FROM developer WHERE pid = ?"

	rows, err := db.Query(query, pid)
	if err != nil {
		return DevList{}, err
	}
	defer rows.Close()

	var dList DevList

	// Extracting product name based on PID
	dList.ProdName, err = ExtractingProdName(pid)
	if err != nil {
		return DevList{}, err
	}

	if rows.Next() {
		var uid int
		var uname string

		err := rows.Scan(&uid)
		if err != nil {
			return DevList{}, errors.Wrap(err, "failed to process a row")
		}

		uname, err = extractingUsername(uid)
		if err != nil {
			return DevList{}, err
		}

		dList.Dev = append(dList.Dev, uname)
	} else {
		dList.Msg = "No issue complete"

		return dList, err
	}

	for rows.Next() {
		var uid int
		var uname string

		err := rows.Scan(&uid)
		if err != nil {
			return DevList{}, errors.Wrap(err, "failed to process a row")
		}

		uname, err = extractingUsername(uid)
		if err != nil {
			return DevList{}, err
		}

		dList.Dev = append(dList.Dev, uname)
	}
	
	err = rows.Err()
	if err != nil {
		return DevList{}, errors.Wrap(err, "failed after processing rows")
	}

	dList.Msg = ""
	return dList, nil
}

// Just a temp function
func CheckEmpty(uid int) error {
	db, err := setup()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "SELECT pid FROM developer WHERE uid = ?"

	var pid int
	
	err = db.QueryRow(query, uid).Scan(&pid)
	if err != nil {
		return err
	}

	return nil
}
