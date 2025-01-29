package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"online-shop/controllers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"online-shop/models"

	"github.com/tebeka/selenium"
)

// unit with getProducts()
func TestGetProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %s", err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "name", "price", "category"}).
		AddRow(1, "Product1", 10.0, "Category1").
		AddRow(2, "Product2", 20.0, "Category2")

	// Expect a query and return mock rows
	mock.ExpectQuery("SELECT id, name, price, category FROM products").WillReturnRows(mockRows)

	// Create a request to test the handler
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()

	// Call the handler with the mock database
	handler := controllers.GetProducts(db)
	handler.ServeHTTP(rec, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []models.Product
	err = json.Unmarshal(rec.Body.Bytes(), &products)
	assert.NoError(t, err)

	// Expected products
	expectedProducts := []models.Product{
		{ID: 1, Name: "Product1", Price: 10.0, Category: "Category1"},
		{ID: 2, Name: "Product2", Price: 20.0, Category: "Category2"},
	}
	assert.Equal(t, expectedProducts, products)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

const (
	seleniumPath = "D:\\program_files\\chromedriver\\chromedriver"
	port         = 8080
)

// end to end with catalog
func TestCatalogPage(t *testing.T) {
	// Start Selenium WebDriver
	service, err := selenium.NewChromeDriverService(seleniumPath, port)
	if err != nil {
		log.Fatalf("Error starting ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Connect to WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		t.Fatalf("Failed to open session: %v", err)
	}
	defer driver.Quit()

	// Open the catalog page
	err = driver.Get("http://localhost:8080/static/catalog.html")
	assert.NoError(t, err)

	// Check navigation links
	navLinks, err := driver.FindElements(selenium.ByCSSSelector, "nav a")
	assert.NoError(t, err)
	assert.Len(t, navLinks, 6)

	// Filter products
	categoryInput, err := driver.FindElement(selenium.ByID, "category")
	assert.NoError(t, err)
	categoryInput.SendKeys("Electronics")

	submitButton, err := driver.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	assert.NoError(t, err)
	submitButton.Click()

	// Wait for results
	driver.Wait(func(wd selenium.WebDriver) (bool, error) {
		products, err := wd.FindElements(selenium.ByCSSSelector, "#products div")
		return len(products) > 0, err
	})

	// Pagination test
	nextPage, err := driver.FindElement(selenium.ByID, "nextPage")
	assert.NoError(t, err)
	nextPage.Click()

	// Verify URL changed
	url, err := driver.CurrentURL()
	assert.NoError(t, err)
	assert.Contains(t, url, "page=2")
}

type Response struct {
	Products    []Product `json:"products"`
	PrevPage    int       `json:"prevPage"`
	NextPage    int       `json:"nextPage"`
	CurrentPage int       `json:"currentPage"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

// integration test for products handler
func TestGetCatalog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Define test cases
	tests := []struct {
		name        string
		queryParams string
		mockRows    func()
		expected    Response
	}{
		{
			name:        "Page 1, Default Sorting",
			queryParams: "?page=1",
			mockRows: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "price", "category", "description", "image"}).
					AddRow(1, "Laptop", 999.99, "Electronics", "A powerful laptop", "laptop.jpg").
					AddRow(2, "Phone", 499.99, "Electronics", "A new smartphone", "phone.jpg")

				mock.ExpectQuery(`SELECT id, name, price, category, description, image FROM products LIMIT 10 OFFSET 0`).
					WillReturnRows(rows)
			},
			expected: Response{
				Products: []Product{
					{ID: 1, Name: "Laptop", Price: 999.99, Category: "Electronics", Description: "A powerful laptop", Image: "laptop.jpg"},
					{ID: 2, Name: "Phone", Price: 499.99, Category: "Electronics", Description: "A new smartphone", Image: "phone.jpg"},
				},
				PrevPage:    0,
				NextPage:    2,
				CurrentPage: 1,
			},
		},
		{
			name:        "Page 2, Category Filter",
			queryParams: "?page=2&category=Electronics",
			mockRows: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "price", "category", "description", "image"}).
					AddRow(3, "Tablet", 299.99, "Electronics", "A lightweight tablet", "tablet.jpg")

				mock.ExpectQuery(`SELECT id, name, price, category, description, image FROM products WHERE category = \$1 LIMIT 10 OFFSET 10`).
					WithArgs("Electronics").
					WillReturnRows(rows)
			},
			expected: Response{
				Products: []Product{
					{ID: 3, Name: "Tablet", Price: 299.99, Category: "Electronics", Description: "A lightweight tablet", Image: "tablet.jpg"},
				},
				PrevPage:    1,
				NextPage:    0,
				CurrentPage: 2,
			},
		},
		{
			name:        "Sorting by Price Descending",
			queryParams: "?sort_by=price&order=desc",
			mockRows: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "price", "category", "description", "image"}).
					AddRow(4, "TV", 1200.00, "Electronics", "A 55-inch Smart TV", "tv.jpg").
					AddRow(5, "Headphones", 150.00, "Electronics", "Noise-canceling headphones", "headphones.jpg")

				mock.ExpectQuery(`SELECT id, name, price, category, description, image FROM products ORDER BY price desc LIMIT 10 OFFSET 0`).
					WillReturnRows(rows)
			},
			expected: Response{
				Products: []Product{
					{ID: 4, Name: "TV", Price: 1200.00, Category: "Electronics", Description: "A 55-inch Smart TV", Image: "tv.jpg"},
					{ID: 5, Name: "Headphones", Price: 150.00, Category: "Electronics", Description: "Noise-canceling headphones", Image: "headphones.jpg"},
				},
				PrevPage:    0,
				NextPage:    2,
				CurrentPage: 1,
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockRows()

			// Create a request
			req, err := http.NewRequest("GET", "/catalog"+tc.queryParams, nil)
			assert.NoError(t, err)

			// Mock response writer
			rr := httptest.NewRecorder()
			handler := controllers.GetCatalog(db)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
			var response Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, response)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
