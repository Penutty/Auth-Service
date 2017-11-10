// License to go here

// +build ignore

// This package serves as a definiton of the authservice's API Endpoints.
// Go Generate uses this package to create unit-tests and request body definitions.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/penutty/apidef"
	"net/http"
)

var (
	path      = flag.String("path", "", "Used to identify API URI.")
	method    = flag.String("method", "", "Used to identify API method.")
	genTests  = flag.Bool("tests", false, "Generate tests.")
	genStruct = flag.Bool("struct", false, "Generate struct.")

	ErrorInvalidEndpointParameters = errors.New("-path AND -method flags must be true.")
	ErrorInvalidGenerateParameters = errors.New("-test XOR -struct flags must be true.")
	ErrorApiEndPointDNE            = errors.New("API Endpoint is not defined.")
)

func main() {
	flag.Parse()
	switch {
	case *path == "" || *method == "":
		panic(ErrorInvalidEndpointParameters)
	case (*genTests && *genStruct) || (!*genTests && !*genStruct):
		panic(ErrorInvalidGenerateParameters)
	}

	var e *apidef.EndPoint
	switch {
	case *path == "/user" && *method == http.MethodPost:
		e = apidef.NewEndpoint([]byte("/user"), http.MethodPost)
		_ = e.NewField("userID", "string").NewValidField("alphanumeric").NewValidField("length", "6", "64").PassWith("user1234").FailWith("123")
		_ = e.NewField("email", "string").NewValidField("email").NewValidField("length", "6", "128").PassWith("useremail@email.com").FailWith("notanemail")
		_ = e.NewField("password", "string").NewValidField("alphanumeric").NewValidField("length", "8", "64").PassWith("userpassword").FailWith("fail")

	case *path == "/auth" && *method == http.MethodPost:
		e = apidef.NewEndpoint([]byte("/auth"), http.MethodPost)
		_ = e.NewField("userID", "string").NewValidField("alphanumeric").NewValidField("length", "6", "64").PassWith("user1234").FailWith("123")
		_ = e.NewField("userID", "string").NewValidField("alphanumeric").NewValidField("length", "8", "64").PassWith("userpassword").FailWith("fail")

	default:
		panic(ErrorApiEndPointDNE)
	}

	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	switch {
	case *genTests:
		e.Tests(buf)
	case *genStruct:
		e.Struct(buf)
	}
	fmt.Printf("%s", buf)
}
