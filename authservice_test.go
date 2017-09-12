package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"testutil"
)

func TestMain(m *testing.M) {
	// Create http.Client or http.Transport if necessary

	os.Exit(m.Run())
}

func Test_APICall_createUser_UserDoesNotExist(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":    {"TJP"},
			"Email":     {"TJP@email.com"},
			"FirstName": {"James"},
			"LastName":  {"Perry"},
			"Password":  {"password"}})
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
	_, _ = testutil.DeleteUser("TJP")
}

func Test_APICall_createUser_UserAlreadyExists(t *testing.T) {
	// createUser
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":    {"TJP"},
			"Email":     {"TJP@email.com"},
			"FirstName": {"James"},
			"LastName":  {"Perry"},
			"Password":  {"password"}})
	if err != nil {
		t.Error(err)
	}

	// attempt to createUser that already exists
	resp, err = http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":    {"TJP"},
			"Email":     {"TJP@email.com"},
			"FirstName": {"James"},
			"LastName":  {"Perry"},
			"Password":  {"password"}})
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

	// Cleanup
	_, _ = testutil.DeleteUser("TJP")
}

func Test_APICall_createUser_MissingCredentials(t *testing.T) {
	resp, err := http.PostForm("http://localhost:8080/user/TJP/create",
		url.Values{
			"UserID":    {"TJP"},
			"FirstName": {"James"},
			"LastName":  {"Perry"},
			"Password":  {"Password"}})
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
			"UserID":    {"TJP"},
			"Email":     {"TJP@email.com"},
			"FirstName": {"James"},
			"LastName":  {"Perry"},
			"Password":  {"password"},
			"Extra":     {"Extra"}})
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
			"Email":     {"TJP@email.com"},
			"FirstName": {"James"},
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
			"UserID":   {"user_100"},
			"Password": {"password_100"}})
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
			"UserID":   {"user_200"},
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
