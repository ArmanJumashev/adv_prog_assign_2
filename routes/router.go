package routes

import (
	"database/sql"
	"online-shop/controllers"
	"online-shop/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Product Routes
	router.HandleFunc("/products", controllers.GetProducts(db)).Methods("GET")
	router.HandleFunc("/catalog", controllers.GetCatalog(db)).Methods("GET")
	router.HandleFunc("/product", controllers.GetProductById(db)).Methods("GET")

	// User Routes
	//     router.HandleFunc("/users", controllers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/register", controllers.Register(db)).Methods("POST")
	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	router.HandleFunc("/profile", controllers.GetUserProfile(db)).Methods("GET")
	router.HandleFunc("/update-user", controllers.UpdateUserProfile(db)).Methods("POST")
	router.HandleFunc("/orders", controllers.GetUserOrders(db)).Methods("GET")
	//     router.HandleFunc("/support", controllers.SendSupportMessage()).Methods("POST")
	router.HandleFunc("/support", controllers.SendEmail)
	router.HandleFunc("/confirm", controllers.SendEmail)

	router.HandleFunc("/confirm-email", controllers.ConfirmEmailHandler)

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.AuthToken)

	protectedRouter.HandleFunc("/order", controllers.OrderHandler(db)).Methods("POST")

	return router
}
