package main

import (
	"log"
	"os"

	"github.com/esponges/initial-setup/internal/seeds"
)

// to seed run this command
// go run cmd/spanner/main.go dmlwrite projects/${PROJECT_ID}/instances/test-instance/databases/example-db
// source: https://cloud.google.com/spanner/docs/getting-started/go

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", os.Args[0], " <database>")
	}

	var db string
	if os.Args[1] == "dmlwrite" {
		if len(os.Args) != 3 {
			log.Fatal("Usage: ", os.Args[0], " dmlwrite <database>")
		}
		db = os.Args[2]
	} else {
		db = os.Args[1]
	}

	err := seeds.Write(os.Stdout, db)
	if err != nil {
		log.Fatal(err)
	}
}
