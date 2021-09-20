package main

import (
	"fmt"
	"os"

	"github.com/bazaglia/warehouse/pkg/importing"
	"github.com/bazaglia/warehouse/pkg/storage/postgres"
)

func main() {
	connURI := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"),
	)

	dir, err := os.Getwd()
	if err != nil {
		panic("Cannot get current directory")
	}

	repository, err := postgres.NewStorage(connURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	importer := importing.NewService(repository)
	path := dir + "/pkg/importing/samples/"

	switch os.Args[1] {
	case "product", "products":
		file, err := os.Open(path + "products.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file: %v\n", err)
		}

		if err = importer.ImportProducts(file); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import file: %v\n", err)
		}
	case "article", "articles", "inventory":
		file, err := os.Open(path + "inventory.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file: %v\n", err)
		}

		if err = importer.ImportArticles(file); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import file: %v\n", err)
		}
	default:
		fmt.Println("Import parameters should be: product,article")
		os.Exit(1)
	}
}
