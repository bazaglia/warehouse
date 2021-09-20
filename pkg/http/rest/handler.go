package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bazaglia/warehouse/pkg/importing"
	"github.com/bazaglia/warehouse/pkg/listing"
	"github.com/bazaglia/warehouse/pkg/selling"
	"github.com/julienschmidt/httprouter"
)

func Handler(l listing.Service, s selling.Service, i importing.Service) http.Handler {
	router := httprouter.New()

	router.GET("/products", getProducts(l))
	router.POST("/products/:id/sell", sellProduct(s))
	router.POST("/products", upload(i.ImportProducts))
	router.POST("/articles", upload(i.ImportArticles))

	return router
}

// getProducts returns a handler for GET /products requests
func getProducts(s listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		list, err := s.ListProducts()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to list products", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

// getProducts returns a handler for GET /products requests
func sellProduct(s selling.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := p.ByName("id")
		amount := 1

		if err := s.SellProduct(id, amount); err != nil {
			fmt.Println(err)
			switch err {
			case selling.ErrNotFound:
				http.Error(w, "The product tried to be sold does not exist", http.StatusBadRequest)
				return
			default:
				http.Error(w, "Unexpected error", http.StatusInternalServerError)
				return
			}
		}

		json.NewEncoder(w).Encode(Response{
			Message: "product successfully sold",
			Success: true,
		})
	}
}

// upload returns a handler that passes an input stream to the provided import function
func upload(importFn func(io.Reader) error) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving the file: %v\n", err)
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if err = importFn(file); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import file: %v\n", err)
			http.Error(w, "Error importing file content", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(Response{
			Message: "file importing was successful",
			Success: true,
		})
	}
}
