// Package user is dedicated to reading and writing user data in Auth-Db.
// The User and AuthCredentials resource definitions and methods are below.
package user

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"log"
	"os"
)

var (
	connStr = os.Getenv("DatabaseConnStr")
	driver  = "mssql"

	ErrorEmailParameterInvalid    = errors.New("User.email must be a valid email address.")
	ErrorUserIDParameterInvalid   = errors.New("AuthCredentials.userID must be a valid userID.")
	ErrorPasswordParameterInvalid = errors.New("AuthCredentials.password must be a valid password.")
	ErrorUserRowNotCreated        = errors.New("Create failed to create one row in the user.Users table.")
)

// MomentDB returns a connection to the SQLSRV Moment-Db database.
func MomentDB() *sql.DB {
	momentDb, err := sql.Open(driver, connStr)
	if err != nil {
		panic(err)
	}
	return momentDb
}

// User references a unique user.Users row in the Moment-Db database.
type User struct {
	userID   string
	email    string
	password string
	err      error
}

// NewUser is a constructor of the User struct.
func NewUser(userID, email, password string) (u *User) {
	u = new(User)
	u.setUserID(userID)
	u.setUserEmail(email)
	u.setPassword(password)
	return
}

// setUserEmail sets User.email if email is valid.
func (u *User) setUserEmail(email string) {
	if u.err != nil {
		return
	}
	if len(email) < 8 {
		u.err = ErrorEmailParameterInvalid
		return
	}
	u.email = email
}

// setUserID sets User.userID if userID is valid.
func (u *User) setUserID(userID string) {
	if u.err != nil {
		return
	}
	if err := checkUserID(userID); err != nil {
		u.err = err
		return
	}
	u.userID = userID
}

// checkUserID returns an error if userID in invalid.
func checkUserID(userID string) error {
	if len(userID) < 6 {
		return ErrorUserIDParameterInvalid
	}
	return nil
}

// setPassword sets User.password if password is valid.
func (u *User) setPassword(password string) {
	if u.err != nil {
		return
	}
	if len(password) < 8 {
		u.err = ErrorPasswordParameterInvalid
		return
	}
	u.password = password
}

// Create inserts a new row into the user.Users table in db.
func Create(u *User, db sq.BaseRunner) {
	if u.err != nil {
		return
	}

	insert := sq.Insert("[auth].[Users]").Columns("[UserID]", "[Email]", "[Password]").Values(u.userID, u.email, u.password)
	res, err := insert.RunWith(db).Exec()
	if err != nil {
		log.Print(err)
		u.err = err
		return
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
		u.err = err
		return
	}
	if cnt != 1 {
		log.Print(ErrorUserRowNotCreated)
		u.err = ErrorUserRowNotCreated
	}
	return
}

// Fetch selects a row from the user.Users table in db.
func Fetch(userID string, db sq.BaseRunner) (u *User) {
	if checkUserID(userID) != nil {
		return
	}

	users := sq.Select("[UserID], [Email], [Password]").From("[auth].[Users]")
	user := users.Where(sq.Eq{"[UserID]": userID})

	u = new(User)
	row := user.RunWith(db).QueryRow()
	err := row.Scan(&u.userID, &u.email, &u.password)
	if err != nil {
		log.Print(err)
		u.err = err
	}
	return
}

// Err returns the the error status of a User instance.
func (u *User) Err() error {
	return u.err
}
