package models

import (
	"fmt"
	"time"
)

func genereteID() string {
	return fmt.Sprintf("%d", time.Now().Nanosecond())
}

type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Definition string    `json:"definition"`
	Price      float64   `json:"price"`
	Image      string    `json:"image"`
	CreateDir  time.Time `json:"create_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewProduct(name, definition string, price float64, image string) *Product {
	return &Product{
		ID:         genereteID(),
		Name:       name,
		Definition: definition,
		Price:      price,
		Image:      image,
		CreateDir:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
