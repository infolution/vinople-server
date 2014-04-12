package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"bytes"
	"encoding/json"
)

var echoTests = []struct {
	in string
	out string
}{
	{
		in: "dog",
		out: "ogday",
	},
	{
		in: "cat",
		out: "atcay",
	},
	{
		in: "hat",
		out: "athay",
	},
}

var signupTests = []struct {
	name string
	wineryName string
	email string
	password string
	success bool
}{
	{
		name: "Mislav Kašner",
		wineryName: "OPG Kašner",
		email: "mislav@emgd.hr",
		password: "soad6ff",
		success: true,
	},
	{
		name: "Matija Kašner",
		wineryName: "OPG Kašner",
		password: "soad6ff",
		success: false,
	},

}


func TestEchoHandler(t *testing.T) {
	request,_ := http.NewRequest("GET","http://localhost:8080/echo/?callback=myfunc", nil);
	response := httptest.NewRecorder()
	echoHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, "\"success\":true") {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "success", body)
	}

//	for i, test := range translateTests {
//		actual := Translate(test.in)
//		assert.Equal(t, test.out, actual, "Test %d", i)
//	}
}

func TestSignupHandler(t *testing.T) {
//	request,_ := http.NewRequest("GET","http://localhost:8080/echo/?callback=myfunc", nil);
	var err error
	vinople,err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v",err)
	}
	var buffer bytes.Buffer
	for i, test := range signupTests {
		response := httptest.NewRecorder()
		buffer.WriteString("http://localhost:8080/echo/?")
		if test.name != "" {
		buffer.WriteString("name=")
		buffer.WriteString(test.name)
		}
		if test.wineryName != "" {
		buffer.WriteString("&wineryName=")
		buffer.WriteString(test.wineryName)
		}
		if test.email != "" {
		buffer.WriteString("&email=")
		buffer.WriteString(test.email)
		}
		if test.password != "" {
		buffer.WriteString("&password=")
		buffer.WriteString(test.password)
		}
		request,_ := http.NewRequest("GET",buffer.String(), nil);
		buffer.Reset()
		signupHandler(response, request)
		var actual Response
		err := json.Unmarshal(response.Body.Bytes(), &actual)
		t.Logf("Signup Response %d: %v", i, actual.Data)
		if (err != nil) {
			t.Fatalf("Error %v generating response: %v", err, response.Body.String())
		}
		assert.Equal(t, test.success, actual.Success, "Test %d", i)




	}
}

