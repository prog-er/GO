package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var redisClient *redis.Client
var db *sql.DB

func init() {
	
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,                
	})

	
	var err error
	db, err = sql.Open("postgres", "postgresql://username:password@localhost:5432/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	productID := vars["id"]

	
	cachedData, err := redisClient.Get(r.Context(), productID).Result()
	if err == nil {
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	
	row := db.QueryRow("SELECT id, name, description, price FROM products WHERE id = $1", productID)
	var product Product
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	
	productJSON, _ := json.Marshal(product)

	
	err = redisClient.Set(r.Context(), productID, productJSON, time.Hour).Err()
	if err != nil {
		log.Println("Failed to cache product data in Redis:", err)
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.Write(productJSON)
}

func main() {
	r := mux.NewRouter()

	
	r.HandleFunc("/products/{id}", getProductByIDHandler).Methods("GET")

	
	fmt.Println("Server listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}


type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
