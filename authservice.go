package main

import (
	"encoding/json"
	"errors"
	valid "github.com/asaskevich/govalidator"
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

	valid.SetFieldsRequiredByDefault(false)
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

	switch r.Method {
	case "POST":
		err := a.postUser(r)
		genErrorHandler(w, err)
	default:
		Error.Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (a *app) authHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		token, err := a.postAuth(r)
		genErrorHandler(w, err)

		w.Header().Set("jwt", token)
		w.WriteHeader(http.StatusOK)

	default:
		Error.Println(ErrorMethodNotImplemented)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (a *app) postUser(r *http.Request) error {
	type body struct {
		UserID   string `valid: "alpha, length(6|64)"`
		Email    string `valid: "email, length(8|128)"`
		Password string `valid: "alphanumeric, length(6|64)"`
	}
	b := new(body)
	if err := json.NewDecoder(r.Body).Decode(b); err != nil {
		return err
	}

	u := a.c.NewUser(b.UserID, b.Email, b.Password)
	a.c.Create(u, user.MomentDB())
	log.Println(a.c.Err())
	return a.c.Err()
}

var ErrorInvalidPass = errors.New("Form value \"Password\" is invalid.")

func (a *app) postAuth(r *http.Request) (string, error) {
	type body struct {
		UserID   string `valid: alpha, length(6|64)"`
		Password string `valid: "alphanumeric, length(6|64)"`
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

func genErrorHandler(w http.ResponseWriter, err error) {
	switch err {
	case nil:
		return
	default:
		Error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
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
