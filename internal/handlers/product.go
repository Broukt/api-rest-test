package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Broukt/api-rest-test/internal/database"
	"github.com/Broukt/api-rest-test/internal/models"
)

// RegisterProductRoutes registers all product related endpoints using net/http.
func RegisterProductRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", getProducts)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("GET /products/{id}", getProductByID)
	mux.HandleFunc("PUT /products/{id}", updateProductByID)
	mux.HandleFunc("DELETE /products/{id}", deleteProductByID)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "product not found"})
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	var products []models.Product
	if err := database.DB.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, products)
}

type createProductInput struct {
	Name     *string  `json:"name"`
	SKU      *string  `json:"sku"`
	Price    *float64 `json:"price"`
	Stock    *int     `json:"stock"`
	IsActive *bool    `json:"isActive"`
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input createProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	if input.Name == nil || input.SKU == nil || input.Price == nil || input.Stock == nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing required fields"})
		return
	}

	p := models.Product{
		Name:  *input.Name,
		SKU:   *input.SKU,
		Price: *input.Price,
		Stock: *input.Stock,
	}

	if input.IsActive != nil && !*input.IsActive {
		now := time.Now()
		p.DeletedAt = &now
	}

	if err := database.DB.Create(&p).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

type updateProductInput struct {
	Name  *string  `json:"name"`
	SKU   *string  `json:"sku"`
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
}

func updateProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "product not found"})
		return
	}

	var input updateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.SKU != nil {
		updates["sku"] = *input.SKU
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Stock != nil {
		updates["stock"] = *input.Stock
	}

	if len(updates) == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "no fields to update"})
		return
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "product not found"})
		return
	}

	now := time.Now()
	if err := database.DB.Model(&product).Update("deleted_at", now).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
