package user

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"testutil"
)

const testingUser = "User_100"
const testingEmail = "Email_100@Email.com"

func Test_get_ByUserID(t *testing.T) {
	user := new(User)
	user.UserID = testingUser
	user.get()

	assert.Equal(t, "Email_100@Email.com", user.Email)

	assert.Equal(t, "FirstName_100", user.FirstName)

	assert.Equal(t, "LastName_100", user.LastName)

}

func Test_get_ByEmail(t *testing.T) {
	user := new(User)
	user.Email = testingEmail
	user.get()

	assert.Equal(t, "user_100", user.UserID)

	assert.Equal(t, "Email_100@Email.com", user.Email)

	assert.Equal(t, "FirstName_100", user.FirstName)

	assert.Equal(t, "LastName_100", user.LastName)

	assert.Equal(t, "Password_100", user.Password)

}

func Test_create(t *testing.T) {
	newUserUserID := "User_test"
	newUser := &User{
		AuthCredentials: AuthCredentials{
			UserID:   newUserUserID,
			Password: "Password_test"},
		Email:     "Email_test@Email.com",
		FirstName: "FirstName_test",
		LastName:  "LastName_test"}

	newUser.create()
	actualUser := new(User)
	actualUser.UserID = newUserUserID
	actualUser.get()

	assert.Equal(t, newUser.UserID, actualUser.UserID)

	assert.Equal(t, newUser.Email, actualUser.Email)

	assert.Equal(t, newUser.FirstName, actualUser.FirstName)

	assert.Equal(t, newUser.LastName, actualUser.LastName)

	assert.Equal(t, newUser.Password, actualUser.Password)

	// Cleanup
	db := openDbConn()
	defer db.Close()

	deleteQuery := `DELETE FROM [User].[Users]
					WHERE UserID = ?`
	if _, err := db.Exec(deleteQuery, newUserUserID); err != nil {
		t.Error(err)
	}
}

func Test_setPassword(t *testing.T) {

	newPassword := "New_Test_Password"
	newPasswordUser := new(User)
	newPasswordUser.UserID = testingUser
	newPasswordUser.Password = newPassword

	newPasswordUser.setPassword()
	actualUser := new(User)
	actualUser.UserID = testingUser
	actualUser.get()

	assert.Equal(t, newPasswordUser.Password, actualUser.Password)

	// Cleanup
	db := openDbConn()
	defer db.Close()

	updateQuery := `UPDATE [User].[Users]
					SET Password = ?
					WHERE UserID = ?`

	if _, err := db.Exec(updateQuery, "Password_100", testingUser); err != nil {
		t.Error(err)
	}
}

func Test_setUserEmail(t *testing.T) {
	updateUser := new(User)
	updateUser.UserID = testingUser
	updateUser.Email = "Email_Updated@Email.com"

	updateUser.setUserEmail()

	actualUser := new(User)
	actualUser.UserID = testingUser
	actualUser.get()

	assert.Equal(t, updateUser.Email, actualUser.Email)

	// Cleanup
	db := openDbConn()
	defer db.Close()

	updateQuery := `UPDATE [User].[Users]
					SET Email = ?
					WHERE UserID = ?`

	if _, err := db.Exec(updateQuery, "Email_100@Email.com", testingUser); err != nil {
		t.Error(err)
	}
}

func Test_setUserFirstName(t *testing.T) {
	updateUser := new(User)
	updateUser.UserID = testingUser
	updateUser.FirstName = "FirstName_Updated"

	updateUser.setUserFirstName()

	actualUser := new(User)
	actualUser.UserID = testingUser
	actualUser.get()

	assert.Equal(t, updateUser.FirstName, actualUser.FirstName)

	// Cleanup
	db := openDbConn()
	defer db.Close()

	updateQuery := `UPDATE [User].[Users]
					SET FirstName = ?
					WHERE UserID = ?`

	if _, err := db.Exec(updateQuery, "FirstName_100", testingUser); err != nil {
		t.Error(err)
	}
}

func Test_setUserLastName(t *testing.T) {
	updateUser := new(User)
	updateUser.UserID = testingUser
	updateUser.LastName = "LastName_Updated"

	updateUser.setUserLastName()

	actualUser := new(User)
	actualUser.UserID = testingUser
	actualUser.get()

	assert.Equal(t, updateUser.LastName, actualUser.LastName)

	// Cleanup
	db := openDbConn()
	defer db.Close()

	updateQuery := `UPDATE [User].[Users]
					SET LastName = ?
					WHERE UserID = ?`

	if _, err := db.Exec(updateQuery, "LastName_100", testingUser); err != nil {
		t.Error(err)
	}
}

func Test_CreateUser(t *testing.T) {
	u := new(User)
	u.UserID = "tjp"
	u.Email = "tjp@email.com"
	u.Password = "Password_test"
	u.FirstName = "FirstName_test"
	u.LastName = "LastName_test"

	if err := CreateUser(u); err != nil {
		t.Error(err)
	}

	u.get()
	assert.Exactly(t, "tjp", u.UserID)
	assert.Exactly(t, "tjp@email.com", u.Email)
	assert.Exactly(t, "Password_test", u.Password)
	assert.Exactly(t, "FirstName_test", u.FirstName)
	assert.Exactly(t, "LastName_test", u.LastName)

	// Cleanup
	testutil.DeleteUser("tjp")

}

func Test_AuthCredentials_validate_Pass(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUser
	aC.Password = "password_100"

	err := aC.validate()
	assert.Empty(t, err)
}

func Test_AuthCredentials_validate_Fail(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUser
	aC.Password = "SomeRandomPasswordThatNoOneHas"

	err := aC.validate()
	assert.NotEmpty(t, err)
}

func Test_AuthUser_Pass(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUser
	aC.Password = "Password_100"

	err := AuthUser(aC)
	assert.Empty(t, err)
}

func Test_AuthUser_Fail(t *testing.T) {
	aC := new(AuthCredentials)
	aC.UserID = testingUser
	aC.Password = "SomeRandomPasswordThatNoOneHas"

	err := AuthUser(aC)
	assert.NotEmpty(t, err)
}
