package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"strconv"
	"strings"
	"testing"
)

var (
	tUser          = "testuser"
	tUserShort     = "user"
	tUserLong      = strings.Repeat("u", 65)
	tUserSpecChars = "testUser!?!??!"

	tEmail              = "<testemail@email.com>"
	tEmailShort         = "<e@a.com>"
	tEmailLong          = "<" + strings.Repeat("u", 128) + "@email.com>"
	tEmailInvalidFormat = "<notanemail>"

	tPassword            = "TestPassword123!"
	tPasswordShort       = "Ac1!"
	tPasswordLong        = strings.Repeat(tPassword, 10)
	tPasswordNoSpecChars = "Abcd1234"
)

func Test_setUserEmail(t *testing.T) {
	type emailErrPair struct {
		email string
		err   error
	}
	testVars := []*emailErrPair{
		&emailErrPair{tEmail, nil},
		&emailErrPair{tEmailShort, ErrorEmailShort},
		&emailErrPair{tEmailLong, ErrorEmailLong},
		&emailErrPair{tEmailInvalidFormat, errors.New("mail: missing @ in addr-spec")},
	}

	for i, v := range testVars {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := CheckEmail(v.email)
			if v.err != nil {
				assert.EqualError(t, err, v.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_setUserID(t *testing.T) {
	type userIDErrPair struct {
		userID string
		err    error
	}
	testVars := []*userIDErrPair{
		&userIDErrPair{tUser, nil},
		&userIDErrPair{tUserShort, ErrorUserIDShort},
		&userIDErrPair{tUserLong, ErrorUserIDLong},
		&userIDErrPair{tUserSpecChars, ErrorUserIDInvalidRunes},
	}

	for i, v := range testVars {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := CheckUserID(v.userID)
			if v.err != nil {
				assert.EqualError(t, err, v.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_setPassword(t *testing.T) {
	type passErrPair struct {
		pass string
		err  error
	}
	testVars := []*passErrPair{
		&passErrPair{tPassword, nil},
		&passErrPair{tPasswordShort, ErrorPasswordShort},
		&passErrPair{tPasswordLong, ErrorPasswordLong},
		&passErrPair{tPasswordNoSpecChars, ErrorPasswordSpecChars},
	}

	for i, v := range testVars {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := CheckPassword(v.pass)
			if v.err != nil {
				assert.EqualError(t, err, v.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_NewUser(t *testing.T) {
	uc := new(UserClient)
	u := uc.NewUser(tUser, tEmail, tPassword)
	assert.Nil(t, u.Err())
}

func Test_Err(t *testing.T) {
	uc := new(UserClient)
	u := uc.NewUser(tUserShort, tEmail, tPassword)
	assert.EqualError(t, u.Err(), ErrorUserIDShort.Error())
}

func Test_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occured when opening a stub database connection. ERROR: %v\n", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO \[auth]\.\[Users] \(\[UserID],\[Email],\[Password]\) VALUES \(\?,\?,\?\)`).
		WithArgs(tUser, tEmail, tPassword).
		WillReturnResult(sqlmock.NewResult(1, 1))

	uc := new(UserClient)
	u := uc.NewUser(tUser, tEmail, tPassword)
	uc.Create(u, db)
	assert.Nil(t, u.Err())

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met. ERROR: %v\n", err)
	}
}

func Test_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occured when opening a stub database connection. ERROR: %v\n", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"UserID", "Email", "Password"}).
		AddRow(tUser, tEmail, tPassword)

	mock.ExpectQuery(`SELECT \[UserID], \[Email], \[Password] FROM \[auth]\.\[Users] WHERE \[UserID] = \?`).
		WithArgs(tUser).
		WillReturnRows(row)

	uc := new(UserClient)
	u := uc.Fetch(tUser, db)
	assert.Nil(t, u.Err())

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met. ERROR: %v\n", err)
	}
}
