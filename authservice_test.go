package main

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Penutty/authservice/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var (
	tUser     = "testuser"
	tEmail    = "<testemail@email.com>"
	tPassword = "TestPassword123!"
)

func NewAuthBody(u, p string) *strings.Reader {
	return strings.NewReader(fmt.Sprintf("{\"UserID\": \"%s\", \"Password\": \"%s\"}", u, p))
}

func NewUserBody(u, e, p string) *strings.Reader {
	return strings.NewReader(fmt.Sprintf("{\"UserID\": \"%s\", \"Email\": \"%s\", \"Password\": \"%s\"}", u, e, p))
}

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

type RequestCodePair struct {
	req  *http.Request
	code int
}

func Test_userHandler(t *testing.T) {
	a := new(app)
	a.c = new(MockUserClient)

	testVars := []*RequestCodePair{
		&RequestCodePair{httptest.NewRequest(http.MethodPost, UserEndpoint, NewUserBody(tUser, tEmail, tPassword)), http.StatusCreated},
		&RequestCodePair{httptest.NewRequest(http.MethodPost, UserEndpoint, strings.NewReader("")), http.StatusBadRequest},
		&RequestCodePair{httptest.NewRequest(http.MethodPost, UserEndpoint, NewUserBody("fail", tEmail, tPassword)), http.StatusBadRequest},
		&RequestCodePair{httptest.NewRequest("METHOD_DNE", UserEndpoint, NewUserBody(tUser, tEmail, tPassword)), http.StatusNotImplemented},
	}

	for i, v := range testVars {
		rec := httptest.NewRecorder()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a.userHandler(rec, v.req)
			assert.Equal(t, v.code, rec.Code)
		})
	}
}

func Test_authHandler(t *testing.T) {
	a := new(app)
	a.c = new(MockUserClient)

	testVars := []*RequestCodePair{
		&RequestCodePair{httptest.NewRequest(http.MethodPost, AuthEndpoint, NewAuthBody(tUser, tPassword)), http.StatusOK},
		&RequestCodePair{httptest.NewRequest(http.MethodPost, AuthEndpoint, strings.NewReader("")), http.StatusBadRequest},
		&RequestCodePair{httptest.NewRequest(http.MethodPost, AuthEndpoint, NewAuthBody("fail", tPassword)), http.StatusBadRequest},
		&RequestCodePair{httptest.NewRequest(http.MethodPost, AuthEndpoint, NewAuthBody(tUser, "fail")), http.StatusBadRequest},
		&RequestCodePair{httptest.NewRequest("METHOD_DNE", AuthEndpoint, NewAuthBody(tUser, tPassword)), http.StatusNotImplemented},
	}

	for i, v := range testVars {
		rec := httptest.NewRecorder()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a.authHandler(rec, v.req)
			assert.Equal(t, v.code, rec.Code)
		})
	}
}

type RequestErrPair struct {
	req *http.Request
	err error
}

func Test_postUser(t *testing.T) {
	a := new(app)
	a.c = new(MockUserClient)

	testVars := []*RequestErrPair{
		&RequestErrPair{httptest.NewRequest(http.MethodPost, UserEndpoint, NewUserBody(tUser, tEmail, tPassword)), nil},
	}

	for i, v := range testVars {
		v.req.Header.Add("Content-Type", "application/json")

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := a.postUser(v.req)
			if v.err != nil {
				assert.EqualError(t, err, v.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_postAuth(t *testing.T) {
	a := new(app)
	a.c = new(MockUserClient)

	testVars := []*RequestErrPair{
		&RequestErrPair{httptest.NewRequest(http.MethodPost, AuthEndpoint, NewAuthBody(tUser, tPassword)), nil},
	}

	for i, v := range testVars {
		v.req.Header.Add("Content-Type", "application/json")

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := a.postAuth(v.req)
			if v.err != nil {
				assert.EqualError(t, err, v.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_generateJwt_pass(t *testing.T) {
	tokenString, err := generateJwt(tUser)
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
		assert.Equal(t, tUser, claims["sub"])
	} else {
		t.Error("token.Claims.(jwt.MapClaims) assertion failed.")
	}

}
