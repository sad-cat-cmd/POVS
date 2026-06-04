package models

import (
	"errors"
	"fmt"
	"strings"
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
	CreateAt   time.Time `json:"create_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func checkName(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	if len(name) > 50 {
		return errors.New("length name > 50")
	}
	return nil
}
func checkDefinition(definition string) (string, error) {
	if definition == "" {
		return "no data", nil
	}
	if len(definition) > 500 {
		return "", errors.New("length difinition > 500")
	}
	return definition, nil
}
func checkPrice(price float64) error {
	if price <= 0 {
		return errors.New("value price <= 0")
	}
	if price > 100000000000 {
		return errors.New("value price > 100000000000")
	}
	return nil
}
func checkImage(image string) (string, error) {
	correctnessExpand := []string{"png", "jpg", "jpeg"}
	if image == "" {
		return "no data", nil
	}
	if len(image) > 255 {
		return "", errors.New("length image > 255")
	}
	complexStr := strings.Split(image, ".")

	if len(complexStr) == 1 || len(complexStr) > 2 {
		return "", errors.New("image contains multiple extensions")
	}
	for _, ch := range complexStr[0] {
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= 'а' && ch <= 'я') ||
			(ch >= 'А' && ch <= 'Я') ||
			(ch >= '0' && ch <= '9') ||
			(ch == '_') ||
			(ch == '-') {
			continue
		} else {
			return "", errors.New("image contains unvalid chars")
		}
	}
	for _, expand := range correctnessExpand {
		if complexStr[1] == expand {
			return image, nil
		}
	}
	return "", errors.New("image contains unvalid extensions")
}
func NewProduct(name string,
	definition string,
	price float64,
	image string) (*Product, error) {
	var errorsMsg string
	var flagError bool
	var newDefinition string
	var newImage string

	err := checkName(name)
	if err != nil {
		flagError = true
		errorsMsg += err.Error() + ";"
	}
	newDefinition, err = checkDefinition(definition)
	if err != nil {
		flagError = true
		errorsMsg += err.Error() + ";"
	}
	err = checkPrice(price)
	if err != nil {
		flagError = true
		errorsMsg += err.Error() + ";"
	}
	newImage, err = checkImage(image)
	if err != nil {
		flagError = true
		errorsMsg += err.Error() + ";"
	}

	if flagError {
		return nil, errors.New(errorsMsg)
	}
	return &Product{
		ID:         genereteID(),
		Name:       name,
		Definition: newDefinition,
		Price:      price,
		Image:      newImage,
		CreateAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}
