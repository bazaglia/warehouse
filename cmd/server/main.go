package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bazaglia/warehouse/pkg/http/rest"
	"github.com/bazaglia/warehouse/pkg/importing"
	"github.com/bazaglia/warehouse/pkg/listing"
	"github.com/bazaglia/warehouse/pkg/selling"
	"github.com/bazaglia/warehouse/pkg/storage/postgres"
)

func main() {
	connURI := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"),
	)

	storage, err := postgres.NewStorage(connURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	lister := listing.NewService(storage)
	seller := selling.NewService(storage)
	importer := importing.NewService(storage)

	router := rest.Handler(lister, seller, importer)

	fmt.Println("The warehouse server is running on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
