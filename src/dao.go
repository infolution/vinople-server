package main

import (
	"text/template"
	"io"
	"bytes"
//	"github.com/satori/go.uuid"

)

//Queries
const (
	readQueryTpl   = "select data from {{.Table}} where clientid = $1 and lastmodified > $2"
	createQueryTpl = "insert into {{.Table}} (id, data, clientid, isdeleted, lastmodified) values ($1,$2,$3,$4,$5)"
	selectByIdTpl  = "select 1 from {{.Table}} where id = $1"
	updateQueryTpl = "update {{.Table}} set data=$1, clientid=$2, isdeleted=$3, lastmodified=$4 where id=$5"
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
	Id string	`json:"id"`
	FullName string	`json:"fullName"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	Email string	`json:"email"`
	Active bool	`json:"active"`

}

//Represents winery which uses this app
type Client struct {
	Id string	`json:"id"`
	WineryName string	`json:"wineryName"`
	TaxNr string	`json:"taxNr"`
	Address string	`json:"address"`
	PostalCode string	`json:"postalCode"`
	City string	`json:"city"`
	Country string	`json:"country"`
	Active bool	`json:"active"`
}

var (
	readQuery  *template.Template
	createQuery *template.Template
)

func init() {
	var err error
	readQuery, err = template.New("readQuery").Parse(readQueryTpl)
	if err != nil { panic(err) }
	createQuery, err = template.New("createQuery").Parse(createQueryTpl)
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

