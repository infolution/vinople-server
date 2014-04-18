package main

import (
	"fmt"
	"net/http"
	"io"
	"net/url"
	"log"
	"encoding/json"
	"time"
	"github.com/satori/go.uuid"
	"github.com/sirsean/go-mailgun/mailgun"
	"code.google.com/p/go.net/websocket"
	"strings"
)

type ResponseModel struct {
	StoreName	string	`json:"storeName"`
	Output	[]string	`json:"output"`
}

type Response struct {
	Success bool    `json:"success"`
	Result	string	`json:"result"`
	Error   string    `json:"error"`
	Models    []ResponseModel    `json:"models"`
	Uid	string	`json:"uid"`
	Action	string	`json:"action"`
}

type SyncRequest struct {
	Action string	`json:action`
	Resources []SyncResource	`json:resources`
	Clientid	string	`json:clientid`
	Uid	string	`json:"uid"`
}
type SyncResource struct {
	Resource 	string	`json:resource`
	StoreName	string	`json:storeName`
	Lastmodified	int64	`json:lastmodified`
	Data	map[string]interface {}	`json:data`
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
	nameparam := q["fullName"]
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
			response.Error = "Error writing client to database: "+err.Error()
			response.Success = false
		}
		//response.Data = string(data)
		loginUserId := uuid.NewV4().String()
		loginUser := &LoginUser{Id: loginUserId, FullName: nameparam[0], Active: true, Username: emailparam[0], Email: emailparam[0], Password: passwordparam[0]}
		data, err = json.Marshal(loginUser)
		if err != nil {
			response.Error = "Error marshaling LoginUser"
		}
		dbquery = createQueryTemplate("login_User")
		_, err = vinople.db.Exec(dbquery, loginUser.Id, string(data), client.Id, false, now)
		if err != nil {
			response.Error = "Error writing loginUser to database: "+err.Error()
			response.Success = false
		}

		sent := email(loginUser.Email,"hr")
		if !sent {
			fmt.Println("Error sending confirmation email")
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	q, _ := url.ParseQuery(query)

	callbackparam := q["callback"]
	clientidparam := q["clientid"]

	usernameparam := q["username"]
	passwordparam := q["password"]
	if callbackparam == nil {

	} else {
		callback := callbackparam[0]
		fmt.Fprintf(w, "%s(", callback)
	}
	response := Response{}
	response.Success = false

	if clientidparam == nil || usernameparam == nil || passwordparam == nil {
		response.Error = "Complete data not provided"
	} else {

		response.Success = true
		//My code start


		dbquery := loginQueryTpl
		rows := vinople.db.QueryRow(dbquery, clientidparam[0], usernameparam[0],passwordparam[0])
		var login string

		if err := rows.Scan(&login); err != nil {
			log.Panic("Error matching username or password: %v", err)
			response.Error = "Username or password error"
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

func email(email string, lang string) bool {
	subject := "Vinople signup confirmation"
	body := `You've registered for Vinople.
Vinople will allow you to manage your winery better and more efficient.

You can always login at the following URL
http://www.vinople.com

We are always making our product better, so if you run into some problems with
implementation or any other issue that's bothering you feel free to contact us at
support@vinople.com.

Thank you,
Mislav Kasner
Vinople founder`
	if lang == "hr" {
	subject = "Vinople potvrda registracije"
	body = `Registrirali ste se za Vinople.
Vinople vam omogućuje da bolje i učinkovitije upravljate svojom vinarijom.

Uvijek se možete prijaviti u aplikaciju preko sljedećeg linka.
http://www.vinople.com

Mi stalno unapređujemo naš proizvod, stoga ako naiđete na problem prilikom korištenja
slobodno nas kontaktirajte na
support@vinople.com.


Hvala,
Mislav Kašner
Vinople`
	}
	message := mailgun.Message{
		FromName:    "Vinople",
		FromAddress: "support@vinople.com",
		ToAddress:   email,
		BCCAddressList: []string{"mislav@infolution.biz"},
		Subject:     subject,
		Body:        body,
	}
	body, err := vinople.mailgunClient.Send(message)
	if err != nil {
		log.Panic("Error sending confirmation e-mail")
	}
	return true
}
func syncHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	q, _ := url.ParseQuery(query)

	callbackparam := q["callback"]
	log.Printf("Request: %v", r)
	decoder := json.NewDecoder(r.Body)
	var syncReq SyncRequest
	err := decoder.Decode(&syncReq)
	if err != nil {
		log.Printf("Error decoding request: %v \n",err)
	}
//	log.Println(syncReq.Resources)
	//Moram proći kroz sve resource
	models := make([]ResponseModel, 0)
	for _, value := range syncReq.Resources {
//		log.Println(value.StoreName)

		model := selectFromTimestamp(value.Lastmodified,syncReq.Clientid, value.StoreName)
		models = append(models, *model)
	}
	response := &Response{true,"success","", models, syncReq.Uid, syncReq.Action}
	b, err := json.Marshal(response)
	if err != nil {
		log.Panic("Unable to marshal Response")
	}

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
//	w.Header().Add("Content-Type", "application/json")

	if callbackparam == nil {

	} else {
		callback := callbackparam[0]
		fmt.Fprintf(w, "%s(", callback)
	}

	fmt.Fprint(w, string(b[:]))

	if callbackparam == nil {

	} else {
		fmt.Fprintf(w, ");\n")
	}

}
//Mijenja Resource name za table name, i provjerava neke posebne slučajeve koji
//se sastoje od dvije riječi
//S ovom funkcijom moram biti siguran da ću dobiti stvarni naziv tablice
func normalizeTableName(name string) string {
	var normalized string
	switch name {
	case "LoginUser" : normalized = "login_user"
	case "ArtiklStanje": normalized = "artikl_stanje"
	case "LotRepromaterijal": normalized = "lot_repromaterijal"
	case "PosudaSekcija": normalized = "posuda_sekcija"
	default : normalized = strings.ToLower(name)
	}
	return normalized
}

func selectFromTimestamp(lastmodified int64, clientid string, name string) *ResponseModel {
	query := readQueryTemplate(normalizeTableName(name))
	rows, err := vinople.db.Query(query, clientid, lastmodified)
	if err != nil {
		log.Panic(err)
	}

	dataArray :=  make([]string,0)
	for rows.Next() {

		var data string
		if err := rows.Scan(&data); err != nil {
			log.Panic(err)
		}
		dataArray = append(dataArray, data)
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}
	model := &ResponseModel{name, dataArray}
	return model
}



//Used only for testing
func socketEchoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

//Used for handling socket
func socketHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])


	m, err := ws.Write(msg[:n])
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}

