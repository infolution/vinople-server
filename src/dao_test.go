package main

import (
	//	"github.com/stretchr/testify/assert"
	"testing"
//	"testing/iotest"
//	"os"
	"bytes"
	"time"
	"github.com/satori/go.uuid"
)



func TestPrepareTemplates(t *testing.T) {
//	w := os.Stdout
	var buffer bytes.Buffer

	prepareTemplates(&buffer)
	t.Log(buffer.String())
}
func TestReadQueryTemplate(t *testing.T) {
	query := readQueryTemplate("artikl")
	t.Log(query)
}

func TestDBWrite(t *testing.T) {
	vinople,err := initVinople(postgresconnTest)
	defer vinople.db.Close()
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v",err)
	}
	query := createQueryTemplate("login_user")
	id := uuid.NewV4()
	now  := time.Now().Unix()
	_, err = vinople.db.Exec(query,id.String(),"{}",nil,false,now)
	if err != nil {
		t.Fatalf("Error writing to database: %v", err)
	}

}

