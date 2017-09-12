// Package resource is dedicated to managing this API's resources.
// The User resource definition and methods are below.
package user

import (
	"database/sql"
	"errors"
	_ "github.com/minus5/gofreetds"
	"log"
)

var UserAlreadyExists = errors.New("User already exists.")
var UserCreateFailed = errors.New("User create failed.")
var ConnStrFailed = errors.New("Unable to connect to SQL Server.")
var UserObjectIsEmpty = errors.New("User is empty.")
var UniquifierIsEmpty = errors.New("User.UserID AND User.Email is empty.")
var UserIDIsEmpty = errors.New("User.UserID AND User.Email is empty.")
var UserEmailIsEmpty = errors.New("User.Email is empty.")
var UserFirstNameIsEmpty = errors.New("User.FirstName is empty.")
var UserLastNameIsEmpty = errors.New("User.LastName is empty.")
var UserPasswordIsEmpty = errors.New("User.Password is empty.")

// User references a unique person represented by a Authentication Credentials and general information.
type User struct {
	AuthCredentials
	Email     string
	FirstName string
	LastName  string
}

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

// Select uses a UserID OR Email to SELECT the corresponding User data
func (u *User) get() error {

	db := openDbConn()
	defer db.Close()

	query := `SELECT UserID,
					 Email, 
				  	 FirstName,
					 LastName, 
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

	return db.QueryRow(query, predicate).Scan(&u.UserID, &u.Email, &u.FirstName, &u.LastName, &u.Password)
}

// create INSERTs a new User.UserID, User.Email, User.FirstName, User.LastName, and User.Password data combination.
func (u *User) create() (sql.Result, error) {

	db := openDbConn()
	defer db.Close()

	query := `INSERT INTO [User].[Users] (UserID, Email, FirstName, LastName, Password)
			  VALUES (?, ?, ?, ?, ?)`

	res, err := db.Exec(query, u.UserID, u.Email, u.FirstName, u.LastName, u.Password)
	return res, err
}

// SetPassword UPDATEs a current User's User.Password
func (u *User) setPassword() {

	db := openDbConn()
	defer db.Close()

	query := `UPDATE [User].[Users]
			  SET Password = ?
			  WHERE UserID = ?`

	if _, err := db.Exec(query, u.Password, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// SetUserEmail UPDATEs the User.Email value associated with User.UserID
func (u *User) setUserEmail() {

	db := openDbConn()
	defer db.Close()

	query := `UPDATE [User].[Users] 
			  SET Email = ?
			  WHERE UserID = ?`

	if _, err := db.Exec(query, u.Email, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// SetUserFirstName UPDATEs the User.FirstName value associated with the User.UserID
func (u *User) setUserFirstName() {

	db := openDbConn()
	defer db.Close()

	query := `UPDATE [User].[Users]
			  SET FirstName = ?
			  WHERE UserID = ?`

	if _, err := db.Exec(query, u.FirstName, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// SetUserLastName UPDATEs the User.LastName value associated with the User.UserID
func (u *User) setUserLastName() {

	db := openDbConn()
	defer db.Close()

	query := `UPDATE [User].[Users]
			  SET LastName = ?
			  WHERE UserID = ?`

	if _, err := db.Exec(query, u.LastName, u.UserID); err != nil {
		log.Fatal(err)
	}
}

// AuthCredentials are a user's unique authentication credentials.
type AuthCredentials struct {
	UserID   string
	Password string
}

func (aC *AuthCredentials) validate() (err error) {
	db := openDbConn()
	defer db.Close()

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

	return db.QueryRow(query, params...).Scan(&aC.UserID)
}

//
// The functions below are utility functions and are only used in this package.
// There functionality has been abstracted and from the above functions for the
// sake of simplicity and readability.

// openDbConn is a wrapper for sql.Open() with logging.
func openDbConn() *sql.DB {
	driver := "mssql"
	connStr := "Server=192.168.1.4:1433;Database=Auth-Db;User Id=Reader;Password=123"

	dbConn, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatal(ConnStrFailed)
	}

	return dbConn
}
