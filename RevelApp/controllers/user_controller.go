package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	m "week2/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	query := "SELECT * FROM users"
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	if name != "" {
		query += " WHERE name='" + name + "'"
	}

	if age != "" {
		if name != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age=" + age
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.UserType, &user.Password, &user.Email); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var u m.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	u.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users (name, age, address, usertype, password, email) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(u.Name, u.Age, u.Address, u.UserType, u.Password, u.Email)

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	u.ID = int(id)

	var response m.UserResponse
	response.Status = 201
	response.Message = fmt.Sprintf("Inserted user with ID = %d", id)
	response.Data = u
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var u m.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE users SET name = ?, age = ?, address = ?, usertype = ?, password = ?, email = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.Name, u.Age, u.Address, u.UserType, u.Password, u.Email, id)

	fmt.Fprintf(w, "Updated user with ID = %s", id)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
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
		fmt.Fprintf(w, "Deleted user with ID = %s", id)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func CreateUserV2(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var u m.UserV2
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&u)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User created")
}

func UpdateUserV2(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var u m.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Model(&m.User{}).Where("id = ?", id).Updates(u)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User updated")
}

func DeleteUserV2(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	vars := mux.Vars(r)
	id := vars["id"]

	result := db.Delete(&m.User{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User deleted")
}

func GetAllUsersV2(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var users []m.User
	result := db.Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetUsersByAgeV2(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	age := r.URL.Query().Get("age")

	var users []m.User
	result := db.Raw("SELECT * FROM users WHERE age = ?", age).Scan(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
