package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"week2/models"

	"github.com/gorilla/mux"
)

type ProductResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    models.Product `json:"data"`
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO products (Name, Price) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(p.Name, p.Price)
	if err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	p.ID = int(id)
	var response ProductResponse
	response.Status = 201
	response.Message = fmt.Sprintf("Inserted product with ID = %d", id)
	response.Data = p
	json.NewEncoder(w).Encode(response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(p.Name, p.Price, id)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	var response ProductResponse
	response.Status = 200
	response.Message = fmt.Sprintf("Updated product with ID = %s", id)
	response.Data = p
	json.NewEncoder(w).Encode(response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	stmt, err := db.Prepare("DELETE FROM transactions WHERE ProductID = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Failed to delete transactions", http.StatusInternalServerError)
		return
	}

	stmt, err = db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	var response ProductResponse
	response.Status = 200
	response.Message = fmt.Sprintf("Deleted product with ID = %s and all related transactions", id)
	json.NewEncoder(w).Encode(response)
}
