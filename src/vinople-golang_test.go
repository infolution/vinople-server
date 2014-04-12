package main

import (
//	"github.com/stretchr/testify/assert"
	"testing"
)

const postgresconnTest = "../postgresconn-test"

func TestInitVinople(t *testing.T) {
	var err error
	vinople,err = initVinople(postgresconnTest)
	if err != nil {
		t.Fatalf("Error when initializing Vinople: %v",err)
	}
	if vinople.db == nil {
		t.Fatalf("Database not initialized")
	}
}

