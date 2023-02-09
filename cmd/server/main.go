package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/evgeniron/API-Validator/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := store.NewInMemoryDB()
	if err != nil {
		return err
	}
	srv := NewServer(db)
	return http.ListenAndServe(":5000", srv)
}
