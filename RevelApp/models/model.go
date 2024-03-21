package models

import "gorm.io/gorm"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	UserType int    `json:"user_type"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type DetailTransactionsResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []DetailTransaction `json:"data"`
}

type DetailTransaction struct {
	ID       int     `json:"id"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserV2 struct {
	gorm.Model
	Name     string
	Age      uint8
	Address  string
	UserType int
	Password *string
	Email    *string
}
