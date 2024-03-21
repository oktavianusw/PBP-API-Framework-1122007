package controllers

import (
	"encoding/json"
	"net/http"
	m "week2/models"
)

func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	name := r.URL.Query()["name"]

	if len(name) == 0 {
		sendErrorResponse(w, "Name parameter is missing")
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE name=?", name[0])

	var user m.User
	if err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.UserType, &user.Password, &user.Email); err != nil {
		sendErrorResponse(w, "User not found")
	} else {
		generateToken(w, user.ID, user.Name, user.UserType)
		sendSuccessResponse(w, "Success")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response m.UserResponse
	response.Status = 200
	response.Message = "Success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
