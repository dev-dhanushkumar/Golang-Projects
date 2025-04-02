package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dev-dhanushkumar/Golang-Projects/product-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductHandler handles HTTP request for product operations
type ProductHandler struct {
	collection *mongo.Collection
}

// NewProductHandler creates a new product handler
func NewProductHandler(collection *mongo.Collection) *ProductHandler {
	return &ProductHandler{collection: collection}
}

// HandleProducts handlers all product endpoints
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.Contains(r.URL.Path, "/products/") && len(r.URL.Path) > len("/products/") {
			h.GetProduct(w, r)
		} else {
			h.getProducts(w, r)
		}
	case http.MethodPost:
		h.CreateProduct(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.deleteProduct(w, r)
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

// GetProducts retrival all products
func (h *ProductHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup options for pagination
	limit := int64(10)
	page := int64(0)

	// Parse query parameter
	if r.URL.Query().Get("limit") != "" {
		limitParam := r.URL.Query().Get("limit")
		if limitInt, err := primitive.ParseDecimal128(limitParam); err == nil {
			limit = int64(limitInt.String()[0])
		}
	}

	if r.URL.Query().Get("page") != "" {
		pageParam := r.URL.Query().Get("page")
		if pageInt, err := primitive.ParseDecimal128(pageParam); err == nil {
			page = int64(pageInt.String()[0])
		}
	}

	// Setup options
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(page * limit)

	cursor, err := h.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Printf("Error finding products: %v", err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products: %v", err)
		http.Error(w, "Failed to decode Products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Error encoding products: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetProduct retrive a single product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract ID from URL Path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	id := parts[len(parts)-1]
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID format", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = h.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Products not found", http.StatusNotFound)
		} else {
			log.Printf("Error finding product: %v", err)
			http.Error(w, "Failed to retrive product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding product: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateProduct creates a new Product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Read and Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		log.Println("Errror in format: ", err)
		http.Error(w, "Invalid Json format", http.StatusBadRequest)
		return
	}

	// Set creation timeout
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.ID = primitive.NewObjectID()

	// Insert Product into dadabase
	result, err := h.collection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Return the created product
	product.ID = result.InsertedID.(primitive.ObjectID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateProduct updates an existing products
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract ID from URL product
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	id := parts[len(parts)-1]
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid Product ID Format!", http.StatusBadRequest)
		return
	}

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var productUpdate models.Product
	if err := json.Unmarshal(body, &productUpdate); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Ensure we don't change the ID
	productUpdate.ID = objectID
	productUpdate.UpdatedAt = time.Now()

	// Update the product
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": productUpdate}

	result, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return the updated product
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(productUpdate); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteProduct deletes a product
func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	id := parts[len(parts)-1]
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID format", http.StatusBadRequest)
		return
	}

	// Delete the product
	result, err := h.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Product successfully deleted"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
