package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"bytes"
	"encoding/json"
	//	"os"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"
	"github.com/satori/go.uuid"
)

var echoTests = []struct {
	in  string
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
	name       string
	wineryName string
	email      string
	password   string
	success    bool
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
var loginTests = []struct {
	clientid string
	username string
	password string
	success  bool
}{
	{
		clientid: "fd80b12f-2e27-4f8f-92a5-b6d00a6bbf5b",
		username: "mislav@emgd.hr",
		password: "soad6ff",
		success: true,
	},
}

func testEchoHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/echo/?callback=myfunc", nil);
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

func testSignupHandler(t *testing.T) {
	//	request,_ := http.NewRequest("GET","http://localhost:8080/echo/?callback=myfunc", nil);
	var err error
	vinople, err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v", err)
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
		request, _ := http.NewRequest("GET", buffer.String(), nil);
		buffer.Reset()
		signupHandler(response, request)
		var actual Response
		err := json.Unmarshal(response.Body.Bytes(), &actual)

		if (err != nil) {
			t.Fatalf("Error %v generating response: %v", err, response.Body.String())
		}
		assert.Equal(t, test.success, actual.Success, "Test %d", i)




	}
}

func testLoginHandler(t *testing.T) {
	//	request,_ := http.NewRequest("GET","http://localhost:8080/echo/?callback=myfunc", nil);
	var err error
	vinople, err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v", err)
	}
	var buffer bytes.Buffer
	for i, test := range loginTests {
		response := httptest.NewRecorder()
		buffer.WriteString("http://localhost:3000/login/?")
		if test.username != "" {
			buffer.WriteString("username=")
			buffer.WriteString(test.username)
		}
		if test.password != "" {
			buffer.WriteString("&password=")
			buffer.WriteString(test.password)
		}
		if test.clientid != "" {
			buffer.WriteString("&clientid=")
			buffer.WriteString(test.clientid)
		}
		request, _ := http.NewRequest("GET", buffer.String(), nil);
		buffer.Reset()
		loginHandler(response, request)
		var actual Response
		err := json.Unmarshal(response.Body.Bytes(), &actual)

		if (err != nil) {
			t.Fatalf("Error %v generating response: %v", err, response.Body.String())
		}
		assert.Equal(t, test.success, actual.Success, "Test %d", i)

	}
}

func testEmail(t *testing.T) {
	var err error
	vinople, err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v", err)
	}
	sent := email("mislav.kasner@gmail.com", "hr")
	assert.Equal(t, true, sent, "Test %d")
	sent = email("mislav.kasner@gmail.com", "en")
	assert.Equal(t, true, sent, "Test %d")
}

func testSyncHandler(t *testing.T) {
	var err error
	vinople, err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v", err)
	}
	resourceArray := make([]SyncResource,0)
//	syncResource := &SyncResource{"artikl", "Artikl", 0, ""}
//	resourceArray = append(resourceArray, *syncResource)
	syncResource := &SyncResource{"loginUser", "LoginUser", 0, nil}
	resourceArray = append(resourceArray, *syncResource)
	syncRequest := &SyncRequest{"sync", resourceArray, "c5fb8dec-ea9f-43b1-a31f-fa4c46fe1a94", uuid.NewV4().String()}

	b, err := json.Marshal(syncRequest)
	if err != nil {
		t.Fatalf("Unable to marshal Sync request")
	}

	request, _ := http.NewRequest("POST", "http://localhost:3000/sync/", bytes.NewReader(b));
	response := httptest.NewRecorder()
	syncHandler(response, request)
	var actual Response
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	if (err != nil) {
		t.Logf("Error %v generating response: %v", err, response.Body.String())
	}
	t.Logf("%+v", response)
	assert.Equal(t, true, actual.Success)
}

//Load models from json to array of models
func loadModelsFromJson(t *testing.T) *map[string]*[]DbModel {
	path := "/home/mikasner/NetBeansProjects/Vinople/vinople-node-yeoman/sampledata/sampledata.js"
	contents, err := ioutil.ReadFile(path) // For read access.
	if err != nil {
		t.Fatal(err)
	}
	allModelsMap := make(map[string]*[]DbModel)
	//	fmt.Println(string(contents))
	//	dataRegExp := regexp.MustCompile(`LoginUser.*\n.*"data" : (.*)`)
	dataRegExp := regexp.MustCompile(`(\w+).*\n.*"data" : (.*)`)
	//	data := regexp.MustCompile(`"data" : (.*)`)
	dataRaw := dataRegExp.FindAllStringSubmatch(string(contents), -1) //Vraća ono šta je data, prvi group match
	//	fmt.Printf("%v", dataRaw)
	for _, i := range dataRaw {

		name := i[1]
		models := make([]DbModel, 10)
		var f []interface{}
		err = json.Unmarshal([]byte(i[2]), &f)
		if err != nil {
			t.Fatal(err)
		}
		time := time.Now().Unix()
		for _, item := range f {
			//		fmt.Printf("%+v\n", item)
			m := item.(map[string]interface{})
			//		fmt.Printf("%s\t%s\t%s\n", m["id"],m["naziv"],m["clientid"])
			//Dodaj syncstatus, isdeleted, lastmodified
			m["isdeleted"] = 0
			m["lastmodified"] = time
			m["syncstatus"] = 0
			id := m["id"].(string)
			data, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			model := &DbModel{id, string(data), m["clientid"].(string), 0, time}
			models = append(models, *model)
//					fmt.Printf("%+v\n\n", model.Data)
		}
		allModelsMap[name] = &models
	}
	return &allModelsMap
}

////Učitava json iz sampledata.js
func testLoadJson(t *testing.T) {

	models := loadModelsFromJson(t)
	if models == nil {
		t.Fatalf("Error loading models")
	}
	for key, _ := range *models {
		fmt.Printf("%s\t\n", key)
	}
}



//
////Spremit ću Json model u bazu
func TestSaveJsonToDB(t *testing.T) {
	var err error
	vinople, err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v", err)
	}

	fmt.Println("Init vinople finished")
	models := loadModelsFromJson(t)
	if models == nil {
		t.Fatalf("Error loading models")
	}
	for key, data := range *models {
		//		fmt.Printf("%s\t\n", key)

		for _, item := range *data {
			if item.Data != "" {
				result := save(normalizeTableName(key), &item)
				fmt.Printf("%v\n\n", result)
			}
		}
	}

}
