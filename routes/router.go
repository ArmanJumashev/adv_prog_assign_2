package routes

import (
    "database/sql"
    "online-shop/controllers"

    "github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
    router := mux.NewRouter()

    // Product Routes
    router.HandleFunc("/products", controllers.GetProducts(db)).Methods("GET")
    router.HandleFunc("/catalog", controllers.GetCatalog(db)).Methods("GET")

    // User Routes
    router.HandleFunc("/users", controllers.GetUsers(db)).Methods("GET")
    router.HandleFunc("/register", controllers.Register(db)).Methods("POST")
    router.HandleFunc("/login", controllers.Login(db)).Methods("POST")

    // Order Routes
    router.HandleFunc("/orders", controllers.GetOrders(db)).Methods("GET")

    return router
}
