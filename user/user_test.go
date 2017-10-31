package user

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

const (
	testuser     = "testuser"
	testemail    = "testemail@email.com"
	testpassword = "testpassword"
)

func Test_setUserEmail(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setUserEmail, testemail, nil)
	testSet(t, u, u.setUserEmail, "dfsdd", ErrorEmailParameterInvalid)
}

func Test_setUserID(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setUserID, testuser, nil)
	testSet(t, u, u.setUserID, "fail", ErrorUserIDParameterInvalid)
}

func Test_setPassword(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setPassword, testpassword, nil)
	testSet(t, u, u.setPassword, "123", ErrorPasswordParameterInvalid)
}

func testSet(t *testing.T, u *User, fn func(string), arg string, expected error) {
	var name string
	if expected == nil {
		name = "Nil Error"
	} else {
		name = expected.Error()
	}

	u.err = nil
	t.Run(name, func(t *testing.T) {
		fn(arg)
		assertError(t, expected, u.Err())
	})
}

func Test_NewUser(t *testing.T) {
	u := NewUser(testuser, testemail, testpassword)
	assertError(t, nil, u.Err())
}

func Test_Err(t *testing.T) {
	u := NewUser("fail", testemail, testpassword)
	assertError(t, ErrorUserIDParameterInvalid, u.Err())
}

func Test_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occured when opening a stub database connection. ERROR: %v\n", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO \[auth]\.\[Users] \(\[UserID],\[Email],\[Password]\) VALUES \(\?,\?,\?\)`).
		WithArgs(testuser, testemail, testpassword).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := NewUser(testuser, testemail, testpassword)
	Create(u, db)
	assertError(t, nil, u.Err())

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
		AddRow(testuser, testemail, testpassword)

	mock.ExpectQuery(`SELECT \[UserID], \[Email], \[Password] FROM \[auth]\.\[Users] WHERE \[UserID] = \?`).
		WithArgs(testuser).
		WillReturnRows(row)

	u := Fetch(testuser, db)
	assertError(t, nil, u.Err())

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met. ERROR: %v\n", err)
	}
}

func assertError(t *testing.T, expected, actual error) {
	if actual != expected {
		t.Fatalf("expected = %v\nactual = %v\n", expected, actual)
	}
}

//func TestMain(m *testing.M) {
//
//	// Create Users for testing
//	u := &User{
//		AuthCredentials: AuthCredentials{
//			UserID:   testingUserID,
//			Password: testingPassword},
//		Email: testingEmail,
//	}
//	_, err := u.create()
//	if err != nil {
//		fmt.Printf("err = %v", err)
//		return
//	}
//
//	// Execute testing functions
//	call := m.Run()
//
//	// Cleanup
//	testutil.DeleteUser(testingUserID)
//
//	os.Exit(call)
//}
//
//func Test_get_ByUserID(t *testing.T) {
//	u := new(User)
//	u.UserID = testingUserID
//	if err := u.get(); err != nil {
//		t.Error(err)
//	}
//
//	assert.Equal(t, testingUserID, u.UserID)
//	assert.Equal(t, testingEmail, u.Email)
//	assert.Equal(t, testingPassword, u.Password)
//}
//
//func Test_get_ByEmail(t *testing.T) {
//	u := new(User)
//	u.Email = testingEmail
//	u.get()
//
//	assert.Equal(t, testingUserID, u.UserID)
//	assert.Equal(t, testingEmail, u.Email)
//	assert.Equal(t, testingPassword, u.Password)
//
//}
//
//func Test_create(t *testing.T) {
//	newUserID := "newUserID_test"
//	newEmail := "newEmail_test@email.com"
//	newPassword := "newPassword_test"
//
//	newU := &User{
//		AuthCredentials: AuthCredentials{
//			UserID:   newUserID,
//			Password: newPassword},
//		Email: newEmail,
//	}
//
//	newU.create()
//	u := new(User)
//	u.UserID = newUserID
//	u.get()
//
//	assert.Equal(t, newU.UserID, u.UserID)
//	assert.Equal(t, newU.Email, u.Email)
//	assert.Equal(t, newU.Password, u.Password)
//
//	testutil.DeleteUser(newUserID)
//}
//
//func Test_setPassword(t *testing.T) {
//
//	newPassword := "newPassword"
//	updatedU := new(User)
//	updatedU.UserID = testingUserID
//	updatedU.Password = newPassword
//
//	updatedU.setPassword()
//	u := new(User)
//	u.UserID = testingUserID
//	u.get()
//
//	assert.Equal(t, newPassword, u.Password)
//
//	// Cleanup
//	updatedU.Password = testingPassword
//	updatedU.setPassword()
//}
//
//func Test_CreateUser(t *testing.T) {
//	UserID := "TJP"
//	Email := "TJP@email.com"
//	Password := "Password"
//	u := new(User)
//	u.UserID = UserID
//	u.Email = Email
//	u.Password = Password
//
//	if err := CreateUser(u); err != nil {
//		t.Error(err)
//	}
//
//	u.get()
//	assert.Exactly(t, UserID, u.UserID)
//	assert.Exactly(t, Email, u.Email)
//	assert.Exactly(t, Password, u.Password)
//
//	// Cleanup
//	testutil.DeleteUser(UserID)
//
//}
//
//func Test_AuthCredentials_validate_Pass(t *testing.T) {
//	aC := new(AuthCredentials)
//	aC.UserID = testingUserID
//	aC.Password = testingPassword
//
//	err := aC.validate()
//	assert.Empty(t, err)
//}
//
//func Test_AuthCredentials_validate_Fail(t *testing.T) {
//	aC := new(AuthCredentials)
//	aC.UserID = testingUserID
//	aC.Password = "SomeRandomPasswordThatNoOneHas"
//
//	err := aC.validate()
//	assert.NotEmpty(t, err)
//}
//
//func Test_AuthUser_Pass(t *testing.T) {
//	aC := new(AuthCredentials)
//	aC.UserID = testingUserID
//	aC.Password = testingPassword
//
//	err := AuthUser(aC)
//	assert.Empty(t, err)
//}
//
//func Test_AuthUser_Fail(t *testing.T) {
//	aC := new(AuthCredentials)
//	aC.UserID = testingUserID
//	aC.Password = "SomeRandomPasswordThatNoOneHas"
//
//	err := AuthUser(aC)
//	assert.NotEmpty(t, err)
//}
