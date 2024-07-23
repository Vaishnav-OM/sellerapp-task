package handlers

import (
	"context"

	"log"

	"github.com/Vaishnav-OM/product-service/db"
	"github.com/Vaishnav-OM/product-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"
	"strconv"
)

var collection *mongo.Collection

// Initialize the MongoDB collection for the service
func InitCollection() {
	collection = db.Client.Database("productdb").Collection("products")
}

// CreateProduct handles the creation of a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product by ID
func UpdateProduct(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	filter := bson.M{"id": id}
	update := bson.M{"$set": updates}
	log.Println(update)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct removes a product by ID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	filter := bson.M{"id": id}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProducts retrieves a list of products with optional filtering and pagination
func GetProducts(c *gin.Context) {
	query := c.Request.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	sortBy := query.Get("sortBy")
	sortOrder := query.Get("sortOrder")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	sort := bson.D{}
	if sortBy != "" {
		order := 1
		if sortOrder == "desc" {
			order = -1
		}
		sort = append(sort, bson.E{Key: sortBy, Value: order})
	}

	filters := bson.M{}
	if query.Get("name") != "" {
		filters["name"] = query.Get("name")
	}
	if query.Get("description") != "" {
		filters["description"] = query.Get("description")
	}
	if query.Get("colour") != "" {
		filters["colour"] = query.Get("colour")
	}
	if query.Get("dimensions") != "" {
		filters["dimensions"] = query.Get("dimensions")
	}
	if query.Get("currencyUnit") != "" {
		filters["currencyUnit"] = query.Get("currencyUnit")
	}
	if query.Get("price") != "" {
		price, err := strconv.ParseFloat(query.Get("price"), 64)
		if err == nil {
			filters["price"] = price
		}
	}

	log.Println("Filters:", filters)

	opts := options.Find().SetSort(sort).SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	cursor, err := collection.Find(context.Background(), filters, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	for cursor.Next(context.Background()) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID retrieves a single product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	filter := bson.M{"id": id}
	var product models.Product
	log.Println(filter)
	err := collection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	log.Println(product)
	c.JSON(http.StatusOK, product)
}
