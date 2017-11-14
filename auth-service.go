package main

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/penutty/authservice/user"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	UserEndpoint = "/user"
	AuthEndpoint = "/auth"
)

var (
	Info  logType = "INFO"
	Warn  logType = "WARN"
	Error logType = "ERROR"

	listenPort = ":8080"
)

type logType string

func logger(t logType) *log.Logger {
	gopath := os.Getenv("GOPATH")
	f, err := os.OpenFile(gopath+"/log/authservice.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return log.New(f, strings.ToUpper(string(t))+": ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

func main() {
	a := new(app)
	a.c = new(user.UserClient)

	http.HandleFunc(UserEndpoint, a.userHandler)
	http.HandleFunc(AuthEndpoint, a.authHandler)

	logger(Error).Fatal(http.ListenAndServe(listenPort, nil))
}

var (
	ErrorMethodNotImplemented = errors.New("Request method is not implemented by API endpoint.")
)

type app struct {
	c user.Client
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := a.postUser(r); err != nil {
			genErrorHandler(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		logger(Error).Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (a *app) authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		token, err := a.postAuth(r)
		if err != nil {
			genErrorHandler(w, err)
			return
		}
		w.Header().Set("jwt", token)
	default:
		logger(Error).Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func genErrorHandler(w http.ResponseWriter, err error) {
	switch err {
	default:
		logger(Error).Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (a *app) postUser(r *http.Request) error {
	type body struct {
		UserID   string
		Email    string
		Password string
	}
	b := new(body)
	if err := json.NewDecoder(r.Body).Decode(b); err != nil {
		return err
	}

	u := a.c.NewUser(b.UserID, b.Email, b.Password)

	if err := a.c.Err(); err != nil {
		return err
	}

	a.c.Create(u, user.MomentDB())
	return a.c.Err()
}

var ErrorInvalidPass = errors.New("Form value \"Password\" is invalid.")

func (a *app) postAuth(r *http.Request) (string, error) {
	type body struct {
		UserID   string
		Password string
	}
	b := new(body)
	if err := json.NewDecoder(r.Body).Decode(b); err != nil {
		return "", err
	}

	u := a.c.Fetch(b.UserID, user.MomentDB())
	if err := a.c.Err(); err != nil {
		return "", err
	}

	if u.Password() != b.Password {
		return "", ErrorInvalidPass
	}

	token, err := generateJwt(b.UserID)
	return token, err
}

// generateJwt uses a requests UserID and a []byte secret to generate a JSON web token.
func generateJwt(UserID string) (string, error) {

	p, err := ioutil.ReadFile("/home/tjp/.ssh/jwt_private.pem")
	if err != nil {
		return "", err
	}

	t := jwt.New(jwt.SigningMethodRS256)
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	claims["iss"] = "Auth-Service"
	claims["sub"] = UserID
	claims["aud"] = "Moment-Service"
	claims["exp"] = time.Now().UTC().AddDate(0, 0, 7).Unix()
	claims["iat"] = time.Now().UTC().Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM(p)
	if err != nil {
		return "", err
	}

	token, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}
