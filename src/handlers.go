package main

import (
	"fmt"
	"net/http"

	"net/url"
	"log"
	"encoding/json"
	"time"
	"github.com/satori/go.uuid"
)

type Response struct {
	Success bool    `json:"success"`
	Error   string    `json:"error"`
	Data   string    `json:"data"`
}

const (
	seed = "vinople"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	q, _ := url.ParseQuery(query)

	callbackparam := q["callback"]
	if callbackparam == nil {

	} else {
		callback := callbackparam[0]
		fmt.Fprintf(w, "%s(", callback)
	}
	response := Response{}
	response.Success = true

	b, err := json.Marshal(response)
	if err != nil {
		log.Panic("Unable to marshal Response")
	}
	fmt.Fprintf(w, string(b[:]))


	//	fmt.Fprintf(w, "\"success\": %t", success)


	if callbackparam == nil {

	} else {
		fmt.Fprintf(w, ");\n")
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	q, _ := url.ParseQuery(query)

	callbackparam := q["callback"]
	nameparam := q["name"]
	wineryNameparam := q["wineryName"]
	emailparam := q["email"]
	passwordparam := q["password"]
	if callbackparam == nil {

	} else {
		callback := callbackparam[0]
		fmt.Fprintf(w, "%s(", callback)
	}
	response := Response{}
	response.Success = false

	if nameparam == nil || wineryNameparam == nil || emailparam == nil || passwordparam == nil {
		response.Error = "Complete data not provided"
	} else {

		response.Success = true
		//My code start
		now := time.Now().Unix()
		//Create client
		clientId := uuid.NewV4().String()
		client := &Client{Id: clientId, WineryName: wineryNameparam[0], Active: true}
		data, err := json.Marshal(client)
		if err != nil {
			response.Error = "Error marshaling client"
		}
		dbquery := createQueryTemplate("client")
		_, err = vinople.db.Exec(dbquery, client.Id, string(data), client.Id, false, now)
		if err != nil {
			log.Panic("Error writing client to database: %v", err)
			response.Error = "Error writing client to database: " + err.Error()
			response.Success = false
		}
		response.Data = string(data)
		loginUserId := uuid.NewV4().String()
		loginUser := &LoginUser{Id: loginUserId, FullName: nameparam[0], Active: true, Username: emailparam[0], Email: emailparam[0], Password: passwordparam[0]}
		data, err = json.Marshal(loginUser)
		if err != nil {
			response.Error = "Error marshaling LoginUser"
		}
		dbquery = createQueryTemplate("login_User")
		_, err = vinople.db.Exec(dbquery, loginUser.Id, string(data), client.Id, false, now)
		if err != nil {
			response.Error = "Error writing loginUser to database: "+ err.Error()
			response.Success = false
		}

		//My code end
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Panic("Unable to marshal Response")
	}
	fmt.Fprintf(w, string(b[:]))


	if callbackparam == nil {

	} else {
		fmt.Fprintf(w, ");\n")
	}
}

func hello(name string) string {

	return "Hello, " + name + "\n"
}

func Translate(in string) string {
	first := in[0:1]
	rest := in[1:]
	return rest + first + "ay"
}

