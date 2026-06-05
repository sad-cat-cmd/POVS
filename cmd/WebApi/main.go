package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sad-cat-cmd/WebApi/internal/config"
	"github.com/sad-cat-cmd/WebApi/internal/handlers"

	//"github.com/sad-cat-cmd/WebApi/internal/services"
	"github.com/sad-cat-cmd/WebApi/internal/services/data"
)

var pathConfigFile string = "config.json"

func main() {
	configuration, err := config.NewConfiguration(pathConfigFile)
	if err != nil {
		fmt.Printf("Error starting server: failed to read config file: %v\n", err)
		return
	}

	service, err := data.NewProductServiceSQLlite(configuration)
	if err != nil {
		fmt.Println("Error starting server: \n\t ", err.Error())
		return
	}
	defer service.CloseDB()
	productHandler := handlers.NewProductHandler(service)

	http.HandleFunc("GET /products", productHandler.GetAllHandler)
	http.HandleFunc("GET /products/", productHandler.SearchHandler)
	http.HandleFunc("POST /products", productHandler.AddHandler)
	http.HandleFunc("PUT /products/", productHandler.EditHandler)
	http.HandleFunc("DELETE /products/", productHandler.RemoveHandler)

	fmt.Println("WebApi starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
