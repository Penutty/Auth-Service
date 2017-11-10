package main

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/dgrijalva/jwt-go"
	"github.com/penutty/authservice/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var (
	tUser          = "testuser"
	tUserShort     = "user"
	tUserLong      = strings.Repeat("user", 20)
	tUserSpecChars = "testUser!?!??!"

	tEmail              = "testemail@email.com"
	tEmailShort         = "e@a.com"
	tEmailLong          = tUserLong + "@email.com"
	tEmailInvalidFormat = "notanemail"

	tPassword            = "TestPassword123!"
	tPasswordShort       = "abc123"
	tPasswordLong        = strings.Repeat("abc123", 10)
	tPasswordNoSpecChars = "abcd1234"
)

type MockUserClient struct {
	err error
}

func (m *MockUserClient) NewUser(UserID, Email, Password string) *user.User {
	uc := new(user.UserClient)
	u := uc.NewUser(UserID, Email, Password)
	m.err = uc.Err()
	return u
}

func (m *MockUserClient) Fetch(u string, db sq.BaseRunner) *user.User {
	uc := new(user.UserClient)
	return uc.NewUser(u, tEmail, tPassword)
}

func (m *MockUserClient) Create(u *user.User, db sq.BaseRunner) {

}

func (m *MockUserClient) Err() error {
	return m.err
}

type RequestErrPair struct {
	req *http.Request
	err error
}

func Test_postUser(t *testing.T) {
	a := new(app)
	a.c = new(MockUserClient)

	newpair := func(user, email, pass string, err error) *RequestErrPair {
		return &RequestErrPair{
			httptest.NewRequest(http.MethodPost, "/user",
				strings.NewReader(`{
						"UserID": "`+user+`",
						"Email": "`+email+`",
						"Password": "`+pass+`"
					}`)),
			err,
		}
	}

	testVars := []*RequestErrPair{
		newpair(tUser, tEmail, tPassword, nil),
		newpair(tUserShort, tEmail, tPassword, user.ErrorUserIDShort),
		newpair(tUserLong, tEmail, tPassword, user.ErrorUserIDLong),
		newpair(tUserSpecChars, tEmail, tPassword, user.ErrorUserIDInvalidRunes),
		newpair(tUser, tEmailShort, tPassword, user.ErrorEmailShort),
		newpair(tUser, tEmailLong, tPassword, user.ErrorEmailLong),
		newpair(tUser, tEmailInvalidFormat, tPassword, nil),
		newpair(tUser, tEmail, tPasswordShort, user.ErrorPasswordShort),
		newpair(tUser, tEmail, tPasswordLong, user.ErrorPasswordLong),
		newpair(tUser, tEmail, tPasswordNoSpecChars, user.ErrorPasswordSpecChars),
	}

	for i, v := range testVars {
		v.req.Header.Add("Content-Type", "application/json")

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := a.postUser(v.req)
			if err != nil {
				assert.EqualError(t, v.err, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_postAuth(t *testing.T) {
}

func Test_generateJwt_pass(t *testing.T) {
	UserID := "tjp"
	tokenString, err := generateJwt(UserID)
	if err != nil {
		t.Error(err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("jwt.Token[\"alg\") = %v instead of RS256")
		}

		p, err := ioutil.ReadFile("/home/tjp/.ssh/jwt_public.pem")
		if err != nil {
			return nil, err
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(p)
		if err != nil {
			return nil, err
		}

		return key, nil
	})
	if err != nil {
		t.Error(err)
	}

	assert.True(t, token.Valid)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, UserID, claims["sub"])
	} else {
		t.Error("token.Claims.(jwt.MapClaims) assertion failed.")
	}

}

//const testingUserID = "User_Test"
//const testingEmail = "Email_Test@Email.com"
//const testingPassword = "Password_Test"
//
//func TestMain(m *testing.M) {
//
//	// Create Users for testing
//	u := &user.User{
//		AuthCredentials: user.AuthCredentials{
//			UserID:   testingUserID,
//			Password: testingPassword},
//		Email: testingEmail,
//	}
//	err := user.CreateUser(u)
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
//func Test_APICall_createUser_UserDoesNotExist(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
//		url.Values{
//			"UserID":   {"TJP"},
//			"Email":    {"TJP@email.com"},
//			"Password": {"password"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 201
//	actual := resp.StatusCode
//	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Body
//	res, err := testutil.ReadJson(resp.Body)
//	assert.Empty(t, res)
//
//	// Cleanup
//	_, err = testutil.DeleteUser("TJP")
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func Test_APICall_createUser_UserAlreadyExists(t *testing.T) {
//
//	// attempt to createUser that already exists
//	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
//		url.Values{
//			"UserID":   {testingUserID},
//			"Email":    {testingEmail},
//			"Password": {testingPassword}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 409
//	actual := resp.StatusCode
//	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Body
//	res, err := testutil.ReadJson(resp.Body)
//	assert.Empty(t, res)
//
//}
//
//func Test_APICall_createUser_MissingCredentials(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
//		url.Values{
//			"UserID":   {"TJP"},
//			"Password": {"Password"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 400
//	actual := resp.StatusCode
//	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Body
//	res, err := testutil.ReadJson(resp.Body)
//	assert.Empty(t, res)
//}
//
//func Test_APICall_createUser_ExtraCredentials(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
//		url.Values{
//			"UserID":   {"TJP"},
//			"Email":    {"TJP@email.com"},
//			"Password": {"password"},
//			"Extra":    {"Extra"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 400
//	actual := resp.StatusCode
//	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Body
//	res, err := testutil.ReadJson(resp.Body)
//	assert.Empty(t, res)
//}
//
//func Test_APICall_createUser_KeyDoesNotExist(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
//		url.Values{
//			"UserID":    {"TJP"},
//			"MadeUpKey": {"Perry"},
//			"Password":  {"password"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 400
//	actual := resp.StatusCode
//	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Body
//	res, err := testutil.ReadJson(resp.Body)
//	assert.Empty(t, res)
//}
//
//func Test_APICall_authUser_Pass(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/auth",
//		url.Values{
//			"UserID":   {testingUserID},
//			"Password": {testingPassword}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 200
//	actual := resp.StatusCode
//	assert.Exactly(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))
//
//	// Header
//	jwt := resp.Header.Get("jwt")
//	assert.NotEmpty(t, jwt)
//}
//
//func Test_APICall_authUser_Fail_InvalidUser(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/auth",
//		url.Values{
//			"UserID":   {"DNE"},
//			"Password": {"password_ThatNoOneHas"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 401
//	actual := resp.StatusCode
//	assert.Exactly(t, expected, actual)
//
//	// Header
//	jwt := resp.Header.Get("jwt")
//	assert.Empty(t, jwt)
//}
//
//func Test_APICall_authUser_Fail_InvalidUserPasswordMatch(t *testing.T) {
//	resp, err := http.PostForm("http://localhost:8080/auth",
//		url.Values{
//			"UserID":   {testingUserID},
//			"Password": {"password_100"}})
//	if err != nil {
//		t.Error(err)
//	}
//	defer resp.Body.Close()
//
//	// StatusCode
//	expected := 401
//	actual := resp.StatusCode
//	assert.Exactly(t, expected, actual)
//
//	// Header
//	jwt := resp.Header.Get("jwt")
//	assert.Empty(t, jwt)
//}
