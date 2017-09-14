package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"testutil"
)

const testingUserID = "User_Test"
const testingEmail = "Email_Test@Email.com"
const testingPassword = "Password_Test"

func TestMain(m *testing.M) {

	// Create Users for testing
	u := &User{
		AuthCredentials: AuthCredentials{
			UserID:   testingUserID,
			Password: testingPassword},
		Email: testingEmail,
	}
	_, err := u.create()
	if err != nil {
		fmt.Printf("err = %v", err)
		return
	}

	// Execute testing functions
	call := m.Run()

	// Cleanup
	testutil.DeleteUser(testingUserID)

	os.Exit(call)
}

func Test_get_ByUserID(t *testing.T) {
	u := new(User)
	u.UserID = testingUserID
	if err := u.get(); err != nil {
		t.Error(err)
	}

	assert.Equal(t, testingUserID, u.UserID)
	assert.Equal(t, testingEmail, u.Email)
	assert.Equal(t, testingPassword, u.Password)
}

func Test_get_ByEmail(t *testing.T) {
	u := new(User)
	u.Email = testingEmail
	u.get()

	assert.Equal(t, testingUserID, u.UserID)
	assert.Equal(t, testingEmail, u.Email)
	assert.Equal(t, testingPassword, u.Password)

}

func Test_create(t *testing.T) {
	newUserID := "newUserID_test"
	newEmail := "newEmail_test@email.com"
	newPassword := "newPassword_test"

	newU := &User{
		AuthCredentials: AuthCredentials{
			UserID:   newUserID,
			Password: newPassword},
		Email: newEmail,
	}

	newU.create()
	u := new(User)
	u.UserID = newUserID
	u.get()

	assert.Equal(t, newU.UserID, u.UserID)
	assert.Equal(t, newU.Email, u.Email)
	assert.Equal(t, newU.Password, u.Password)

	testutil.DeleteUser(newUserID)
}

func Test_setPassword(t *testing.T) {

	newPassword := "newPassword"
	updatedU := new(User)
	updatedU.UserID = testingUserID
	updatedU.Password = newPassword

	updatedU.setPassword()
	u := new(User)
	u.UserID = testingUserID
	u.get()

	assert.Equal(t, newPassword, u.Password)

	// Cleanup
	updatedU.Password = testingPassword
	updatedU.setPassword()
}

func Test_CreateUser(t *testing.T) {
	UserID := "TJP"
	Email := "TJP@email.com"
	Password := "Password"
	u := new(User)
	u.UserID = UserID
	u.Email = Email
	u.Password = Password

	if err := CreateUser(u); err != nil {
		t.Error(err)
	}

	u.get()
	assert.Exactly(t, UserID, u.UserID)
	assert.Exactly(t, Email, u.Email)
	assert.Exactly(t, Password, u.Password)

	// Cleanup
	testutil.DeleteUser(UserID)

}

func Test_AuthCredentials_validate_Pass(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUserID
	aC.Password = testingPassword

	err := aC.validate()
	assert.Empty(t, err)
}

func Test_AuthCredentials_validate_Fail(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUserID
	aC.Password = "SomeRandomPasswordThatNoOneHas"

	err := aC.validate()
	assert.NotEmpty(t, err)
}

func Test_AuthUser_Pass(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUserID
	aC.Password = testingPassword

	err := AuthUser(aC)
	assert.Empty(t, err)
}

func Test_AuthUser_Fail(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUserID
	aC.Password = "SomeRandomPasswordThatNoOneHas"

	err := AuthUser(aC)
	assert.NotEmpty(t, err)
}
