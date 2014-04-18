package main

import (
	"fmt"
	"net/http"
	//	"net/url"
	//	"encoding/csv"
	//	"os"
	//	"io"
		"io/ioutil"
	"log"
	"flag"
//	"errors"
	//	"time"
	//	"strconv"
	//	"crypto/rand"
	//	"math/big"
	_ "github.com/lib/pq"
	"database/sql"
	//	"github.com/go-martini/martini"
	//	"github.com/satori/go.uuid"
	"github.com/sirsean/go-mailgun/mailgun"
//	"code.google.com/p/go.net/websocket"
)

type Vinople struct {

	db *sql.DB
	mailgunClient *mailgun.Client

}

var dbFlag = flag.String("db", "", "Postgres configuration")

func init() {
	// example with short version for long flag
	flag.StringVar(dbFlag, "d", "", "Postgres configuration")
}

var vinople *Vinople



func initVinople(dbconn string) (*Vinople, error) {

	//Mailgun
	mailgunClient := mailgun.NewClient("key-1tbrmvpihgzt2hzac-pguge25zx8m0w4", "vinople.com")

	//Database
	buff, err := ioutil.ReadFile(dbconn)
	if err != nil {
		return nil, err
	}
	connstring := string(buff)
	//fmt.Printf("%s\n", connstring)
	db, err := sql.Open("postgres", connstring)

	if err != nil {
		return nil, err
	}
	v := &Vinople{db: db, mailgunClient: mailgunClient}
	return v, nil
}

func main() {
	fmt.Println("Initializing Vinople Server...")
	flag.Parse()
	var err error
	if *dbFlag == "" {
		log.Fatal("Postgres configuration is missing")
	}
	vinople, err = initVinople(*dbFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting vinople api server")

	http.HandleFunc("/echo/", echoHandler)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/sync/", syncHandler)
//	http.Handle("/socket/echo/", websocket.Handler(socketEchoHandler))
//	http.Handle("/socket/", websocket.Handler(socketHandler))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
