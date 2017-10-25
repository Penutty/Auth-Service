// Package user is dedicated to reading and writing user data in Auth-Db.
// The User and AuthCredentials resource definitions and methods are below.
package user

import (
	"database/sql"
	"errors"
	"github.com/penutty/dba"
	"log"
	"os"
)

var ConnStr = os.Getenv("DatabaseConnStr")

var UserAlreadyExists = errors.New("User already exists.")
var UserCreateFailed = errors.New("User create failed.")
var ConnStrFailed = errors.New("Unable to connect to SQL Server.")
var UserObjectIsEmpty = errors.New("User is empty.")
var UniquifierIsEmpty = errors.New("User.UserID AND User.Email is empty.")
var UserIDIsEmpty = errors.New("User.UserID AND User.Email is empty.")
var UserEmailIsEmpty = errors.New("User.Email is empty.")
var UserPasswordIsEmpty = errors.New("User.Password is empty.")

// CreateUser checks to see if a user exists and creates it if not.
// If the user already exists and err is returned.
func CreateUser(u *User) (err error) {
	if err = u.get(); err != sql.ErrNoRows {
		return UserAlreadyExists
	}

	res, err := u.create()
	if err != nil {
		return UserCreateFailed
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return UserCreateFailed
	}

	return nil
}

// AuthUser checks to see if the AuthCredentials passed identify a user.
// error will be nil on successful authentication.
func AuthUser(aC *AuthCredentials) (err error) {
	if err = aC.validate(); err != nil {
		return err
	}

	return nil
}

//
// The methods and functions below are NOT exported.
// They are utilized by the exported methods.
//

// User references a unique person represented by a Authentication Credentials and general information.
type User struct {
	AuthCredentials
	Email string
}

// Select uses a UserID OR Email to SELECT the corresponding User data
func (u *User) get() error {
	c := dba.OpenConn(ConnStr)
	defer c.Db.Close()

	query := `SELECT UserID,
					 Email, 
					 Password
			  FROM [User].[Users]`
	var predicate string
	if len(u.UserID) > 0 {
		query = query + ` WHERE UserID = ?`
		predicate = u.UserID
	} else if len(u.Email) > 0 {
		query = query + ` WHERE Email = ?`
		predicate = u.Email
	} else {
		return UniquifierIsEmpty
	}

	return c.Db.QueryRow(query, predicate).Scan(&u.UserID, &u.Email, &u.Password)
}

// create INSERTs a new User.UserID, User.Email, and User.Password data combination.
func (u *User) create() (sql.Result, error) {
	c := dba.OpenConn(ConnStr)
	defer c.Db.Close()

	query := `INSERT INTO [User].[Users] (UserID, Email, Password)
			  VALUES (?, ?, ?)`

	res, err := c.Db.Exec(query, u.UserID, u.Email, u.Password)
	return res, err
}

// SetPassword UPDATEs a current User's User.Password
func (u *User) setPassword() {
	c := dba.OpenConn(ConnStr)
	defer c.Db.Close()

	query := `UPDATE [User].[Users]
			  SET Password = ?
			  WHERE UserID = ?`

	if _, err := c.Db.Exec(query, u.Password, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// SetUserEmail UPDATEs the User.Email value associated with User.UserID
func (u *User) setUserEmail() {
	c := dba.OpenConn(ConnStr)
	defer c.Db.Close()

	query := `UPDATE [User].[Users] 
			  SET Email = ?
			  WHERE UserID = ?`

	if _, err := c.Db.Exec(query, u.Email, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// AuthCredentials are a user's unique authentication credentials.
type AuthCredentials struct {
	UserID   string
	Password string
}

func (aC *AuthCredentials) validate() (err error) {
	c := dba.OpenConn(ConnStr)
	defer c.Db.Close()

	query := `SELECT UserID
			  FROM [user].[Users]
			  WHERE Password = ?
			  AND `
	params := []interface{}{aC.Password}

	if aC.UserID != "" {
		query = query + `UserID = ?`
		params = append(params, aC.UserID)
	} else {
		return UniquifierIsEmpty
	}

	return c.Db.QueryRow(query, params...).Scan(&aC.UserID)
}
