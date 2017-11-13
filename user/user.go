// Package user is dedicated to reading and writing user data in Auth-Db.
// The User and AuthCredentials resource definitions and methods are below.
package user

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/minus5/gofreetds"
	"log"
	"net/mail"
	"os"
	"regexp"
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

type Client interface {
	Newer
	Creater
	Fetcher
	Err() error
}
type CreateFetcher interface {
	Creater
	Fetcher
}

type Newer interface {
	NewUser(string, string, string) *User
}

type Creater interface {
	Create(*User, sq.BaseRunner)
}

type Fetcher interface {
	Fetch(string, sq.BaseRunner) *User
}

type UserClient struct {
	err error
}

// NewUser is a constructor of the User struct.
func (uc *UserClient) NewUser(userID, email, password string) (u *User) {
	u = new(User)
	u.setUserID(userID)
	u.setUserEmail(email)
	u.setPassword(password)
	uc.err = u.err
	return
}

// Create inserts a new row into the user.Users table in db.
func (uc *UserClient) Create(u *User, db sq.BaseRunner) {
	if uc.err != nil {
		return
	}

	insert := sq.Insert("[auth].[Users]").Columns("[UserID]", "[Email]", "[Password]").Values(u.userID, u.email, u.password)
	res, err := insert.RunWith(db).Exec()
	if err != nil {
		log.Print(err)
		uc.err = err
		return
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
		uc.err = err
		return
	}
	if cnt != 1 {
		log.Print(ErrorUserRowNotCreated)
		uc.err = ErrorUserRowNotCreated
	}
	return
}

// Fetch selects a row from the user.Users table in db.
func (uc *UserClient) Fetch(userID string, db sq.BaseRunner) (u *User) {
	if uc.err != nil {
		return
	}
	if err := CheckUserID(userID); err != nil {
		uc.err = err
		return
	}

	users := sq.Select("[UserID], [Email], [Password]").From("[auth].[Users]")
	user := users.Where(sq.Eq{"[UserID]": userID})

	u = new(User)
	row := user.RunWith(db).QueryRow()
	err := row.Scan(&u.userID, &u.email, &u.password)
	if err != nil {
		log.Print(err)
		uc.err = err
	}
	return
}

// Err returns the the error status of a User instance.
func (uc *UserClient) Err() error {
	return uc.err
}

// User references a unique user.Users row in the Moment-Db database.
type User struct {
	userID   string
	email    string
	password string
	err      error
}

// setUserEmail sets User.email if email is valid.
func (u *User) setUserEmail(email string) {
	if u.err != nil {
		return
	}
	if err := CheckEmail(email); err != nil {
		u.err = err
		return
	}
	u.email = email
}

var (
	EmailMinLength = 10
	EmailMaxLength = 128

	ErrorEmailShort = errors.New("Email too short.")
	ErrorEmailLong  = errors.New("Email too long.")
)

// CheckEmail returns an error if email is invalid.
func CheckEmail(email string) error {
	switch {
	case len(email) < EmailMinLength:
		return ErrorEmailShort
	case len(email) > EmailMaxLength:
		return ErrorEmailLong
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}

// setUserID sets User.userID if userID is valid.
func (u *User) setUserID(userID string) {
	if u.err != nil {
		return
	}
	if err := CheckUserID(userID); err != nil {
		u.err = err
		return
	}
	u.userID = userID
}

var (
	UserIDMinLength = 8
	UserIDMaxLength = 64

	ErrorUserIDShort        = errors.New("UserID too short.")
	ErrorUserIDLong         = errors.New("UserID too long.")
	ErrorUserIDInvalidRunes = errors.New("UserID may only consist of numbers and letters.")
)

// CheckUserID returns an error if userID is invalid.
func CheckUserID(userID string) error {
	r, err := regexp.Compile(`^[a-zA-Z0-9]+$`)
	if err != nil {
		return err
	}
	switch {
	case len(userID) < UserIDMinLength:
		return ErrorUserIDShort
	case len(userID) > UserIDMaxLength:
		return ErrorUserIDLong
	case !r.MatchString(userID):
		return ErrorUserIDInvalidRunes
	default:
		return nil
	}
}

// setPassword sets User.password if password is valid.
func (u *User) setPassword(password string) {
	if u.err != nil {
		return
	}
	if err := CheckPassword(password); err != nil {
		u.err = err
		return
	}
	u.password = password
}

var (
	PasswordMinLength = 8
	PasswordMaxLength = 64

	ErrorPasswordShort     = errors.New("Password too short.")
	ErrorPasswordLong      = errors.New("Password too long.")
	ErrorPasswordLowerCase = errors.New("Password does not contain a lowercase letter.")
	ErrorPasswordUpperCase = errors.New("Password does not contain an uppercase letter.")
	ErrorPasswordNumber    = errors.New("Password does not contain a number.")
	ErrorPasswordSpecChars = errors.New("Password does not contain a special character.")
)

// CheckPassword returns an error if password is invalid.
func CheckPassword(password string) error {
	var (
		err error
		ll  *regexp.Regexp
		ul  *regexp.Regexp
		num *regexp.Regexp
		sc  *regexp.Regexp
	)

	if ll, err = regexp.Compile(`[a-z]+`); err != nil {
		return err
	}
	if ul, err = regexp.Compile(`[A-Z]+`); err != nil {
		return err
	}
	if num, err = regexp.Compile(`[0-9]+`); err != nil {
		return err
	}
	if sc, err = regexp.Compile(`[^a-zA-Z0-9]+`); err != nil {
		return err
	}

	switch {
	case len(password) < PasswordMinLength:
		return ErrorPasswordShort
	case len(password) > PasswordMaxLength:
		return ErrorPasswordLong
	case !ll.MatchString(password):
		return ErrorPasswordLowerCase
	case !ul.MatchString(password):
		return ErrorPasswordUpperCase
	case !num.MatchString(password):
		return ErrorPasswordNumber
	case !sc.MatchString(password):
		return ErrorPasswordSpecChars
	}
	return nil
}

func (u *User) Password() (p string) {
	if u.err != nil {
		return
	}
	p = u.password
	return
}

func (u *User) Err() error {
	return u.err
}
