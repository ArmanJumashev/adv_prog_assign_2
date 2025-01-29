package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"online-shop/models"
	"strconv"
)

func GetCatalog(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		category := r.URL.Query().Get("category")
		sortBy := r.URL.Query().Get("sort_by")
		order := r.URL.Query().Get("order")

		// Конвертация номера страницы
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}

		// Базовый SQL-запрос
		query := `SELECT id, name, price, category, description, image FROM products`
		var params []interface{}

		// Фильтрация по категории
		if category != "" {
			query += " WHERE category = $1"
			params = append(params, category)
		}

		// Сортировка
		if sortBy != "" && order != "" {
			query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
		}

		// Пагинация
		limit := 10
		offset := (pageInt - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

		// Выполнение запроса
		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying products: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Парсинг результатов
		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Description, &p.Image); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning product: %v", err), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		// Проверка на следующую страницу
		hasNextPage := len(products) == limit

		// Ответ JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"products": products,
			"prevPage": func() int {
				if pageInt > 1 {
					return pageInt - 1
				} else {
					return 0
				}
			}(),
			"nextPage": func() int {
				if hasNextPage {
					return pageInt + 1
				} else {
					return 0
				}
			}(),
			"currentPage": pageInt,
		})
	}
}

func GetProductById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productIdStr := r.URL.Query().Get("id")
		if productIdStr == "" {
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}

		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		var product models.Product
		err = db.QueryRow("SELECT id, name, description, price, category, image FROM products WHERE id = $1", productId).
			Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Image)

		if err != nil {
			if err.Error() == "no rows in result set" {
				http.Error(w, "Product not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving product", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}
func GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price, category FROM products")
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category); err != nil {
				http.Error(w, "Failed to parse product", http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}
