package main

import (
	"fmt"
	"log"
	"net/http"
	"week2/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db_latihan_pbp.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.Authenticate(controllers.DeleteUser, 1)).Methods("DELETE")

	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")
	// router.HandleFunc("/logout", controllers.Logout).Methods("POST")
	router.HandleFunc("/checkUserLogin", controllers.CheckUserLogin).Methods("GET")

	// V1 Endpoints
	// router.HandleFunc("/v1/users", controllers.GetAllUsers).Methods("GET")
	// router.HandleFunc("/v1/users", controllers.CreateUser).Methods("POST")
	// router.HandleFunc("/v1/users/{id}", controllers.UpdateUser).Methods("PUT")
	// router.HandleFunc("/v1/users/{id}", controllers.DeleteUser).Methods("DELETE")

	// 2 Endpoints
	// router.HandleFunc("/v2/users", controllers.GetAllUsersV2).Methods("GET")
	// router.HandleFunc("/v2/users", controllers.CreateUserV2).Methods("POST")
	// router.HandleFunc("/v2/users/{id}", controllers.UpdateUserV2).Methods("PUT")
	// router.HandleFunc("/v2/users/{id}", controllers.DeleteUserV2).Methods("DELETE")
	// router.HandleFunc("/v2/users/age", controllers.GetUsersByAgeV2).Methods("GET")

	// Product
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	// Transaction
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", controllers.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id}", controllers.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controllers.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/transactions/user/{userID}", controllers.GetDetailUsersTransactions).Methods("GET")
	http.Handle("/", router)
	fmt.Println("Connected to port 1234")
	log.Println("Connected to port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))
}
