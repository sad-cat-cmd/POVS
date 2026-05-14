package main

import (
	"fmt"
	"log"
	"net/http"
)

func main ()
{
	fmt.Println("WebApi starting on port :8080");
	log.Fatal(http.ListenAndServe(":8080", nil));
}