package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sad-cat-cmd/WebApi/internal/interfaces"
	"github.com/sad-cat-cmd/WebApi/internal/models"
)

type ProductHandler struct {
	service interfaces.IProductService
}

func NewProductHandler(serv interfaces.IProductService) *ProductHandler {
	return &ProductHandler{service: serv}
}

// POST /products
func (h *ProductHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	var req_product models.Product

	err := json.NewDecoder(r.Body).Decode(&req_product)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	new_product, errCreateProduct := models.NewProduct(
		req_product.Name,
		req_product.Definition,
		req_product.Price,
		req_product.Image,
	)
	if errCreateProduct != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	created, err := h.service.Add(new_product)
	if err != nil {
		if err.Error() == "Error write data" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// DELETE /products/{id}
func (h *ProductHandler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	removed, err := h.service.Remove(id)

	if err != nil {
		if err.Error() == "Object doesn't exist" {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
		if err.Error() == "Error write data" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(removed)
}

// PUT /products/{id}
func (h *ProductHandler) EditHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product.ID = id
	newProduct, errNewProduct := models.NewProduct(product.Name,
		product.Definition,
		product.Price,
		product.Image)
	if errNewProduct != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	newProduct.ID = id
	updated, err := h.service.Edit(newProduct)
	if err != nil {
		if err.Error() == "object doesn't exist" {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
		if err.Error() == "Error write data" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// GET /products/{id}
func (h *ProductHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	product, err := h.service.Search(id)
	if err != nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GET /products
func (h *ProductHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
