package main

import (
	"text/template"
	"io"
	"bytes"
	//	"github.com/satori/go.uuid"
	"log"

)

//Queries
const (
	readQueryTpl        = "select data from {{.Table}} where clientid = $1 and lastmodified > $2"
	createQueryTpl      = "insert into {{.Table}} (id, data, clientid, isdeleted, lastmodified) values ($1,$2,$3,$4,$5)"
	selectByIdTpl       = "select 1 as exist from {{.Table}} where id = $1"
	updateQueryTpl      = "update {{.Table}} set data=$1, clientid=$2, isdeleted=$3, lastmodified=$4 where id=$5"
	loginQueryTpl       = "select  1 from login_user WHERE clientid=$1 AND data->>'username'  = $2 and  data->>'password'=$3 and data->>'active'='true'"
	selectByClientIdTpl = "select * from {{.Table}} where clientid = $1"
)

type Table struct {
	Table string
}


type DAO struct {
	readQuery   string
	createQuery string
	selectById  string
	updateQuery string
}

//Represents user that can login to application
type LoginUser struct {
	Id       string    `json:"id"`
	FullName string    `json:"fullName"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Active   bool    `json:"active"`

}

//Represents winery which uses this app
type Client struct {
	Id         string    `json:"id"`
	WineryName string    `json:"wineryName"`
	TaxNr      string    `json:"taxNr"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postalCode"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	Active     bool    `json:"active"`
}

type DbModel struct {
	Id              string    `json:"id"`
	Data           string    `json:"data"`
	Clientid        string    `json:"clientid"`
	Isdeleted       int    `json:"isdeleted"`
	Lastmodified    int64    `json:"lastmodified"`
}

var (
	readQuery  *template.Template
	createQuery *template.Template
	updateQuery *template.Template
	selectByClientId *template.Template
	selectById *template.Template
)

func init() {
	var err error
	readQuery, err = template.New("readQuery").Parse(readQueryTpl)
	if err != nil { panic(err) }
	createQuery, err = template.New("createQuery").Parse(createQueryTpl)
	if err != nil { panic(err) }
	updateQuery, err = template.New("updateQuery").Parse(updateQueryTpl)
	if err != nil { panic(err) }
	selectByClientId, err = template.New("selectByClientId").Parse(selectByClientIdTpl)
	if err != nil { panic(err) }
	selectById, err = template.New("selectById").Parse(selectByIdTpl)
	if err != nil { panic(err) }
}

//Kreira template iz querya kako bi mogao zamijeniti tablicu
func prepareTemplates(w io.Writer) {
	tmpl, err := template.New("readQuery").Parse(readQueryTpl)
	if err != nil { panic(err) }

	err = tmpl.Execute(w, &Table{"artikl"})
	if err != nil { panic(err) }
}

//Kreira popunjeni template za query, tj. query koji možemo izvršiti na bazi
//podataka
func readQueryTemplate(tablename string) string {
	var buffer bytes.Buffer
	table := Table{tablename}

	err := readQuery.Execute(&buffer, &table)
	if err != nil { panic(err) }
	return buffer.String()
}


func createQueryTemplate(tablename string) string {
	var buffer bytes.Buffer
	table := Table{tablename}

	err := createQuery.Execute(&buffer, &table)
	if err != nil { panic(err) }
	return buffer.String()
}

func updateQueryTemplate(tablename string) string {
	var buffer bytes.Buffer
	table := Table{tablename}

	err := updateQuery.Execute(&buffer, &table)
	if err != nil { panic(err) }
	return buffer.String()
}
func selectByIdTemplate(tablename string) string {
	var buffer bytes.Buffer
	table := Table{tablename}

	err := selectById.Execute(&buffer, &table)
	if err != nil { panic(err) }
	return buffer.String()
}

func save(tablename string, model *DbModel) int {
	//Prvo provjeravam da li postoji u bazi
	query1 := selectByIdTemplate(tablename)
	rows := vinople.db.QueryRow(query1, model.Id)
	var exist int

	if err := rows.Scan(&exist); err != nil {
//		log.Panic("Error receiving record by id: %v", err)
		exist = -1
	}
	if exist > 0 {    //Ako record postoji radim update
		dbquery := updateQueryTemplate(tablename)
		_, err := vinople.db.Exec(dbquery, model.Data, model.Clientid, model.Isdeleted, model.Lastmodified,model.Id)
		if err != nil {
			log.Panic("Error updating %s to database: %v", tablename, err)
			return 0
		}

	} else {  //Ako ne postoji radim insert
		dbquery := createQueryTemplate(tablename)
		_, err := vinople.db.Exec(dbquery, model.Id, model.Data, model.Clientid, model.Isdeleted, model.Lastmodified)
		if err != nil {
			log.Panic("Error inserting %s to database: %v", tablename, err)
			return 0
		}
	}
	return 1

}

