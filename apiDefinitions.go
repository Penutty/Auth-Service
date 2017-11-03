// Licensing will go here.

// +build ignore

// apiDefinitions.go holds API endpoint definitions and generates tests and structs.
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var (
	path       = flag.String("path", "", "Pass endpoint path to generate for.")
	method     = flag.String("method", "", "Pass endpoint method to generate for.")
	genTests   = flag.Bool("genTests", false, "Use flag to generate tests for endpoint passed as value")
	genStructs = flag.Bool("getStructs", false, "Use flag to generate struct for API definitions defined in main")

	ErrorInvalidEndpointParameters = errors.New("A path and method must be passed via the -path and -method flags")
	ErrorInvalidGenerateParameters = errors.New("The flag genTests OR genStructs must be set. Both may not be set.")
	ErrorApiEndPointDNE            = errors.New("The API endpoint as specified by -path and -method is not defined in main.")
	ErrorInvalidPath               = errors.New("path must be a valid string.")
	ErrorInvalidMethod             = errors.New("method must be a valid string.")
	ErrorInvalidName               = errors.New("name must be a valid string.")
	ErrorInvalidType               = errors.New("type must be a vaild string.")
)

func main() {
	flag.Parse()
	switch {
	case *path == "" || *method == "":
		panic(ErrorInvalidEndpointParameters)
	case (*genTests && *genStructs) || (!*genTests && !*genStructs):
		panic(ErrorInvalidGenerateParameters)
	}

	var endPoints map[endPointKey]*endPoint

	e := NewEndpoint("/user", http.MethodPost)
	_ = e.NewField("userID", "string").NewValidField("alpha").NewValidFieldParams("length", 6, 64)
	_ = e.NewField("email", "string").NewValidField("email").NewValidFieldParams("length", 8, 128)
	_ = e.NewField("password", "string").NewValidField("alphanumeric").NewValidFieldParams("length", 8, 64)
	endPoints[endPointKey{path: e.path, method: e.method}] = e

	e = NewEndpoint("/auth", http.MethodPost)
	_ = e.NewField("userID", "string").NewValidField("alpha").NewValidFieldParams("length", 6, 64)
	_ = e.NewField("password", "string").NewValidField("alphanumeric").NewValidFieldParams("length", 8, 64)
	endPoints[endPointKey{path: e.path, method: e.method}] = e

	ep, ok := endPoints[endPointKey{path: *path, method: *method}]
	if !ok {
		panic(ErrorApiEndPointDNE)
	}

	switch {
	case *genTests:
		ep.GenerateTests()
	case *genStructs:
		ep.GenerateStruct()
	}
}

type endPointKey struct {
	path   string
	method string
}

type endPoint struct {
	endPointKey
	fields []*field
}

func NewEndpoint(path string, method string) *endPoint {
	switch {
	case len(path) <= 0:
		panic(ErrorInvalidPath)
	case len(method) <= 0:
		panic(ErrorInvalidMethod)
	}

	return &endPoint{
		endPointKey: endPointKey{
			path:   path,
			method: method,
		},
	}
}

func (e *endPoint) GenerateTests() {
	fmt.Printf("Generating tests...\n")
}

func (e *endPoint) GenerateStruct() {
	fmt.Printf("type body struct {\n")
	for _, f := range e.fields {
		fmt.Print(f)
	}
	fmt.Printf("}\n")

}

type field struct {
	name       string
	Type       string
	validators []*validField
}

func (e *endPoint) NewField(name string, Type string) *field {
	switch {
	case len(name) <= 0:
		panic(ErrorInvalidName)
	case len(Type) <= 0:
		panic(ErrorInvalidType)
	}

	f := &field{
		name: name,
		Type: Type,
	}
	e.fields = append(e.fields, f)
	return f
}

func (f field) String() {
	fmt.Printf("\t%v %v `valid: \"%v\"`\n", f.name, f.Type, f.validators)
}

type validField struct {
	name string
	min  uint64
	max  uint64
}

func (v validField) String() {
	if v.min == 0 && v.max == 0 {
		fmt.Print(v.name)
	} else {
		fmt.Print(v.name + "(" + strconv.Itoa(int(v.min)) + "|" + strconv.Itoa(int(v.max)) + ")")
	}
}

func (f *field) NewValidField(name string) *field {
	if len(name) <= 0 {
		panic(ErrorInvalidName)
	}

	vf := &validField{
		name: name,
	}
	f.validators = append(f.validators, vf)
	return f
}

func (f *field) NewValidFieldParams(name string, min uint64, max uint64) *field {
	if len(name) <= 0 {
		panic(ErrorInvalidName)
	}

	vf := &validField{
		name: name,
		min:  min,
		max:  max,
	}
	f.validators = append(f.validators, vf)
	return f
}
