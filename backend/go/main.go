package main

import (
	"log"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
	"github.com/ebarrosjr/igreja-aberta/backend/go/router"
)

func main() {
	if err := database.Initialize("igreja_aberta.db"); err != nil {
		log.Fatal(err)
	}

	if err := database.MigrateAll(); err != nil {
		log.Fatal(err)
	}

	if err := database.SeedInitialData(); err != nil {
		log.Fatal(err)
	}

	router.InitializeRouter()
}
