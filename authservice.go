package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/penutty/authservice/user"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	listenPort = ":8080"
)

func init() {
	Logger := func(logType string) *log.Logger {
		file := "/home/tjp/go/log/moment.txt"
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		l := log.New(f, strings.ToUpper(logType)+": ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
		return l
	}

	Info = Logger("info")
	Warn = Logger("warn")
	Error = Logger("error")
}
func main() {
	a := new(app)
	a.c = new(user.UserClient)

	http.HandleFunc("/user", a.userHandler)
	http.HandleFunc("/auth", a.authHandler)

	Error.Fatal(http.ListenAndServe(listenPort, nil))
}

var (
	ErrorMethodNotImplemented = errors.New("Request method is not implemented by API endpoint.")
)

type app struct {
	c user.Client
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		Error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		expected := []string{
			"UserID",
			"Email",
			"Password",
		}
		if err := ValidateForm(r.Form, expected); err != nil {
			Error.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		switch err := a.postUser(r); err {
		case nil:
			break
		default:
			Error.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	default:
		Error.Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (a *app) postUser(r *http.Request) error {
	u := a.c.NewUser(r.FormValue("UserID"), r.FormValue("Email"), r.FormValue("Password"))
	a.c.Create(u, user.MomentDB())
	return a.c.Err()
}

func (a *app) authHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		Error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	switch r.Method {
	case http.MethodPost:
		expected := []string{
			"UserID",
			"Password",
		}
		if err := ValidateForm(r.Form, expected); err != nil {
			Error.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		token, err := a.postAuth(r)
		switch err {
		case nil:
			break
		default:
			Error.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.Header().Set("jwt", token)
		w.WriteHeader(http.StatusOK)

	default:
		Error.Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

var ErrorInvalidPass = errors.New("Form value \"Password\" is invalid.")

func (a *app) postAuth(r *http.Request) (string, error) {
	u := a.c.Fetch(r.FormValue("UserID"), user.MomentDB())
	if err := a.c.Err(); err != nil {
		return "", err
	}

	if u.Password() != r.FormValue("Password") {
		return "", ErrorInvalidPass
	}

	token, err := generateJwt(r.FormValue("UserID"))
	return token, err
}

var (
	ErrorIncorrectNumberOfFields = errors.New("API endpoint expected more fields.")
	ErrorFieldMissing            = errors.New("Request is missing a form field that is required by this API endpoint.")
)

func ValidateForm(f url.Values, expected []string) error {
	if len(f) != len(expected) {
		return ErrorIncorrectNumberOfFields
	}
	var j int
	for i, _ := range f {
		if i != expected[j] {
			return ErrorFieldMissing
		}
		j++
	}
	return nil
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
