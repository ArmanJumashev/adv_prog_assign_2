package data_loader


import (
    "database/sql" // Импортируем пакет database/sql для работы с БД

	"encoding/json"
	"fmt"
	"log"
	"github.com/go-resty/resty/v2"
)

type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

func FetchAndInsertProducts(db *sql.DB) {
	//Fake Store api
	client := resty.New()
	resp, err := client.R().
		SetQueryParam("limit", "40").
		Get("https://fakestoreapi.com/products")

	if err != nil {
		log.Fatalf("Error fetching products: %v", err)
	}

	var products []Product
	if err := json.Unmarshal(resp.Body(), &products); err != nil {
		log.Fatalf("Error unmarshalling products: %v", err)
	}

	for _, product := range products {
		_, err := db.Exec(`
			INSERT INTO products (id, title, price, description, category, image)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			product.ID, product.Title, product.Price, product.Description, product.Category, product.Image)

		if err != nil {
			log.Printf("Error inserting product %d: %v", product.ID, err)
		}
	}

	fmt.Println("Products successfully added to the database.")
}