package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	m "week2/models"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var login m.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user m.User
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", login.Email).Scan(&user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if user.Password != login.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode("User logged in")
}

// func UserLogin(w http.ResponseWriter, r *http.Request) {
//     db := connect()
//     defer db.Close()

//     var login m.Login
//     err := json.NewDecoder(r.Body).Decode(&login)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     var user m.User
//     err = db.QueryRow("SELECT password FROM users WHERE email = ?", login.Email).Scan(&user.Password)
//     if err != nil {
//         if err == sql.ErrNoRows {
//             http.Error(w, "User not found", http.StatusNotFound)
//         } else {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//         return
//     }

//     err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
//     if err != nil {
//         http.Error(w, "Invalid password", http.StatusUnauthorized)
//         return
//     }

//     json.NewEncoder(w).Encode("User logged in")
// }
