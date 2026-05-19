package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sad-cat-cmd/WebApi/internal/models"
	"github.com/sad-cat-cmd/WebApi/internal/services"
)

type ProductHandler struct {
	service services.IProductService
}

func NewProductHandler(serv services.IProductService) *ProductHandler {
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

	new_product := models.NewProduct(
		req_product.Name,
		req_product.Definition,
		req_product.Price,
		req_product.Image,
	)

	created, err := h.service.Add(new_product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusNotFound)
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

	updated, err := h.service.Edit(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// GET /products/{id}
func (h *ProductHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product, err := h.service.Search(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GET /products
func (h *ProductHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
