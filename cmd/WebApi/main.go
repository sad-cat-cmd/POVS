package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sad-cat-cmd/WebApi/cmd/WebApi/handlers"
	"github.com/sad-cat-cmd/WebApi/internal/services"
)

func main() {
	services := services.NewProductServices()
	productHandler := handlers.NewProductHandler(services)

	http.HandleFunc("GET /products", productHandler.GetAllHandler)
	http.HandleFunc("GET /products/", productHandler.SearchHandler)
	http.HandleFunc("POST /products", productHandler.AddHandler)
	http.HandleFunc("PUT /products/", productHandler.EditHandler)
	http.HandleFunc("DELETE /products/", productHandler.RemoveHandler)

	fmt.Println("WebApi starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
