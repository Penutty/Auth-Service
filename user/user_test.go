package user

import (
	"database/sql"
	"testing"
)

func Test_setUserEmail(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setUserEmail, "testuser@email.com", nil)
	testSet(t, u, u.setUserEmail, "dfsdd", ErrorEmailParameterInvalid)
}

func Test_setUserID(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setUserID, "testUser", nil)
	testSet(t, u, u.setUserID, "fail", ErrorUserIDParameterInvalid)
}

func Test_setPassword(t *testing.T) {
	u := new(User)
	testSet(t, u, u.setPassword, "Testpassword", nil)
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
		assert(t, u.Err(), expected)
	})
}

func Test_NewUser(t *testing.T) {
	u := NewUser("testuser", "testemail@email.com", "testpassword")
	assert(t, u.Err(), nil)
}

func Test_Err(t *testing.T) {
	u := NewUser("fail", "testemail@email.com", "testpassword")
	assert(t, u.Err(), ErrorUserIDParameterInvalid)
}

type MockSqBaseRunner struct {
}

func (m *MockSqBaseRunner) Exec(query string, args ...interface{}) (sql.Result, error) {
	res := MockSqlResult{
		lastInsertID: 1,
		rowsAffected: 1,
	}
	return res, nil
}

func (m *MockSqBaseRunner) Query(query string, args ...interface{}) (*sql.Rows, error) {

}

type MockSqlResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (m MockSqlResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m MockSqlResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func Test_Create(t *testing.T) {
	u := NewUser("testUser", "testemail@email.com", "testpassword")
	db := make(MockSqBaseRunner)

	Create(u, db)
	assert(t, u.Err(), nil)
}

//func Test_Fetch(t *testing.T) {
//	db := new(MockSqBaseRunner)
//
//	u := Fetch("testuser", db)
//	assert(t, u.Err(), nil)
//}

func assert(t *testing.T, expected, actual error) {
	if actual != expected {
		t.Fatalf("actual = %v\nexpected = %v\n", actual, expected)
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
