package main

import (
	"log"

	"github.com/Cypher012/OrganizeNoteAPi/internal/automate"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(0)

	if err := automate.RunAutomation(); err != nil {
		log.Fatal(err)
	}
}
