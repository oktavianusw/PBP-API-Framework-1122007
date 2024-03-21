package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	m "week2/models"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	query := "select * from transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var transaction m.Transaction
	var transactions []m.Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	w.Header().Set("Content Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var t m.Transaction
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the product exists
	var productExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE ID = ?)", t.ProductID).Scan(&productExists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the product doesn't exist, insert a new product with an empty name
	if !productExists {
		_, err = db.Exec("INSERT INTO products (ID, Name, Price) VALUES (?, '', 0)", t.ProductID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	stmt, err := db.Prepare("INSERT INTO transactions (UserID, ProductID, Quantity) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(t.UserID, t.ProductID, t.Quantity)

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	var response m.TransactionsResponse
	response.Status = 201
	response.Message = fmt.Sprintf("Inserted transaction with ID = %d", id)
	response.Data = []m.Transaction{t}
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var t m.Transaction
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE transactions SET UserID = ?, ProductID = ?, Quantity = ? WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(t.UserID, t.ProductID, t.Quantity, id)

	fmt.Fprintf(w, "Updated transaction with ID = %s", id)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	stmt, err := db.Prepare("DELETE FROM transactions WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected != 0 {
		fmt.Fprintf(w, "Deleted transaction with ID = %s", id)
	} else {
		http.Error(w, "Transaction not found", http.StatusNotFound)
	}
}

func GetDetailUsersTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	userID := vars["userID"]

	var query string
	if userID != "all" {
		query = fmt.Sprintf(`
            SELECT t.ID, u.ID, u.Name, u.Age, u.Address, p.ID, p.Name, p.Price, t.Quantity
            FROM transactions t
            JOIN users u ON t.UserID = u.ID
            JOIN products p ON t.ProductID = p.ID
            WHERE u.ID = %s
        `, userID)
	} else if userID == "all" {
		query = `
            SELECT t.ID, u.ID, u.Name, u.Age, u.Address, p.ID, p.Name, p.Price, t.Quantity
            FROM transactions t
            JOIN users u ON t.UserID = u.ID
            JOIN products p ON t.ProductID = p.ID
        `
	}

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []m.DetailTransaction
	for rows.Next() {
		var dt m.DetailTransaction
		err := rows.Scan(&dt.ID, &dt.User.ID, &dt.User.Name, &dt.User.Age, &dt.User.Address, &dt.Product.ID, &dt.Product.Name, &dt.Product.Price, &dt.Quantity)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, dt)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read rows", http.StatusInternalServerError)
		return
	}

	var response m.DetailTransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}
