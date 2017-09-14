package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/penutty/authservice/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"testutil"
)

const testingUserID = "User_Test"
const testingEmail = "Email_Test@Email.com"
const testingPassword = "Password_Test"

func TestMain(m *testing.M) {

	// Create Users for testing
	u := &user.User{
		AuthCredentials: user.AuthCredentials{
			UserID:   testingUserID,
			Password: testingPassword},
		Email: testingEmail,
	}
	err := user.CreateUser(u)
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

func Test_APICall_createUser_UserDoesNotExist(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":   {"TJP"},
			"Email":    {"TJP@email.com"},
			"Password": {"password"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 201
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Body
	res, err := testutil.ReadJson(resp.Body)
	assert.Empty(t, res)

	// Cleanup
	_, err = testutil.DeleteUser("TJP")
	if err != nil {
		t.Error(err)
	}
}

func Test_APICall_createUser_UserAlreadyExists(t *testing.T) {

	// attempt to createUser that already exists
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":   {testingUserID},
			"Email":    {testingEmail},
			"Password": {testingPassword}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 409
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Body
	res, err := testutil.ReadJson(resp.Body)
	assert.Empty(t, res)

}

func Test_APICall_createUser_MissingCredentials(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":   {"TJP"},
			"Password": {"Password"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 400
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Body
	res, err := testutil.ReadJson(resp.Body)
	assert.Empty(t, res)
}

func Test_APICall_createUser_ExtraCredentials(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":   {"TJP"},
			"Email":    {"TJP@email.com"},
			"Password": {"password"},
			"Extra":    {"Extra"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 400
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Body
	res, err := testutil.ReadJson(resp.Body)
	assert.Empty(t, res)
}

func Test_APICall_createUser_KeyDoesNotExist(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":    {"TJP"},
			"MadeUpKey": {"Perry"},
			"Password":  {"password"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 400
	actual := resp.StatusCode
	assert.Equal(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Body
	res, err := testutil.ReadJson(resp.Body)
	assert.Empty(t, res)
}

func Test_APICall_authUser_Pass(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/auth",
		url.Values{
			"UserID":   {testingUserID},
			"Password": {testingPassword}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 200
	actual := resp.StatusCode
	assert.Exactly(t, expected, actual, "Status "+string(expected)+" = "+http.StatusText(expected)+"\nStatus "+string(actual)+" = "+http.StatusText(actual))

	// Header
	jwt := resp.Header.Get("jwt")
	assert.NotEmpty(t, jwt)
}

func Test_APICall_authUser_Fail_InvalidUser(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/auth",
		url.Values{
			"UserID":   {"DNE"},
			"Password": {"password_ThatNoOneHas"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 401
	actual := resp.StatusCode
	assert.Exactly(t, expected, actual)

	// Header
	jwt := resp.Header.Get("jwt")
	assert.Empty(t, jwt)
}

func Test_APICall_authUser_Fail_InvalidUserPasswordMatch(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/auth",
		url.Values{
			"UserID":   {testingUserID},
			"Password": {"password_100"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// StatusCode
	expected := 401
	actual := resp.StatusCode
	assert.Exactly(t, expected, actual)

	// Header
	jwt := resp.Header.Get("jwt")
	assert.Empty(t, jwt)
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
