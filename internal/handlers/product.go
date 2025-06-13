package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Broukt/api-rest-test/internal/database"
	"github.com/Broukt/api-rest-test/internal/models"
)

// RegisterProductRoutes registers all product related endpoints.
func RegisterProductRoutes(r *gin.Engine) {
	g := r.Group("/products")
	g.GET("", getProducts)
	g.GET("/:id", getProductByID)
	g.POST("", createProduct)
	g.PUT("/:id", updateProductByID)
	g.DELETE("/:id", deleteProductByID)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func getProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	var products []models.Product
	if err := database.DB.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

type createProductInput struct {
	Name     string  `json:"name" binding:"required"`
	SKU      string  `json:"sku" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Stock    int     `json:"stock" binding:"required"`
	IsActive *bool   `json:"isActive"`
}

func createProduct(c *gin.Context) {
	var input createProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := models.Product{
		Name:  input.Name,
		SKU:   input.SKU,
		Price: input.Price,
		Stock: input.Stock,
	}

	if input.IsActive != nil && !*input.IsActive {
		now := time.Now()
		p.DeletedAt = &now
	}

	if err := database.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

type updateProductInput struct {
	Name  *string  `json:"name"`
	SKU   *string  `json:"sku"`
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
}

func updateProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	var input updateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	now := time.Now()
	if err := database.DB.Model(&product).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
