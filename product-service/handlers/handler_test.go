package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaishnav-OM/product-service/db"
	"github.com/Vaishnav-OM/product-service/models"
	"github.com/gin-gonic/gin"

	// "go.mongodb.org/mongo-driver/bson"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/products", CreateProduct)
	r.PUT("/products/:id", UpdateProduct)
	r.DELETE("/products/:id", DeleteProduct)
	r.GET("/products", GetProducts)
	r.GET("/products/:id", GetProductByID)
	return r
}

func setupTestDB() {
	db.Client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	InitCollection()
}

func populateDB() {
	products := []models.Product{
		{ID: "1", Name: "Product1", Description: "Description1", Colour: "Red", Dimensions: "10x10x10", CurrencyUnit: "USD", Price: 10.00},
		{ID: "2", Name: "Product2", Description: "Description2", Colour: "Blue", Dimensions: "20x20x20", CurrencyUnit: "USD", Price: 20.00},
		{ID: "3", Name: "Product3", Description: "Description3", Colour: "Green", Dimensions: "30x30x30", CurrencyUnit: "USD", Price: 30.00},
		{ID: "4", Name: "Product4", Description: "Description4", Colour: "Yellow", Dimensions: "40x40x40", CurrencyUnit: "USD", Price: 40.00},
		{ID: "5", Name: "Product5", Description: "Description5", Colour: "Black", Dimensions: "50x50x50", CurrencyUnit: "USD", Price: 50.00},
	}

	for _, product := range products {
		collection.InsertOne(context.Background(), product)
	}
}

func clearDB() {
	collection.Drop(context.Background())
}

func TestHandlers(t *testing.T) {
	setupTestDB()
	r := setupRouter()

	t.Run("CreateProduct", func(t *testing.T) {
		clearDB()
		product := models.Product{
			ID:           "6",
			Name:         "Test Product",
			Description:  "Test Description",
			Colour:       "Red",
			Dimensions:   "10x10x10",
			CurrencyUnit: "USD",
			Price:        99.99,
		}
		jsonValue, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		clearDB()
		populateDB()
		updates := map[string]interface{}{
			"name":        "Updated Product",
			"description": "Updated Description",
		}
		jsonValue, _ := json.Marshal(updates)
		req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonValue))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		clearDB()
		populateDB()
		req, _ := http.NewRequest("DELETE", "/products/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("GetProducts", func(t *testing.T) {
		clearDB()
		populateDB()

		// Test case with pagination
		req, _ := http.NewRequest("GET", "/products?page=1&limit=2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var productsResp []models.Product
		err := json.Unmarshal(resp.Body.Bytes(), &productsResp)
		assert.NoError(t, err)
		assert.Len(t, productsResp, 2) // Should return 2 products as per limit

		// Test case with sorting
		req, _ = http.NewRequest("GET", "/products?sortBy=price&sortOrder=desc", nil)
		resp = httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		err = json.Unmarshal(resp.Body.Bytes(), &productsResp)
		assert.NoError(t, err)
		assert.Equal(t, 50.00, productsResp[0].Price) // Should return product with highest price first

		// Test case with filtering
		req, _ = http.NewRequest("GET", "/products?colour=Red", nil)
		resp = httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		err = json.Unmarshal(resp.Body.Bytes(), &productsResp)
		assert.NoError(t, err)
		assert.Equal(t, "Red", productsResp[0].Colour) // Should return product with Red colour
	})

	t.Run("GetProductByID", func(t *testing.T) {
		clearDB()
		populateDB()
		req, _ := http.NewRequest("GET", "/products/2", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var productResp models.Product
		err := json.Unmarshal(resp.Body.Bytes(), &productResp)
		assert.NoError(t, err)
		assert.Equal(t, "Product2", productResp.Name)
	})
}
