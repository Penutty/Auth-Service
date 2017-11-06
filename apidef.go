package main

import (
	"github.com/penutty/apidef"
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
