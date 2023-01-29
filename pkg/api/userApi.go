package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/saddfox/cup-predictor/pkg/auth"
	"github.com/saddfox/cup-predictor/pkg/db"
	"github.com/saddfox/cup-predictor/pkg/models"
)

type registerResponse struct {
	ID   uint   `json:"Id"`
	Name string `json:"Name"`
}

// creates user and returns user id and name
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = user.ValidateUser()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	newUser := models.User{
		Name:     user.Name,
		Email:    strings.ToLower(user.Email),
		Password: hashedPassword,
		Admin:    false,
	}

	result := db.DB.Create(&newUser)

	if result.Error != nil {
		fmt.Println("duplicate user", result.Error.Error())
		ERROR(w, http.StatusConflict, result.Error)
		return
	}

	response := registerResponse{
		ID:   newUser.ID,
		Name: newUser.Name,
	}

	JSON(w, http.StatusCreated, response)
}

// validates user credentials and returns jwt auth_token
func LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var dbUser models.User
	result := db.DB.First(&dbUser, "email = ?", strings.ToLower(user.Email))
	if result.Error != nil {
		ERROR(w, http.StatusUnauthorized, result.Error)
		return
	}

	if auth.VerifyPassword(dbUser.Password, user.Password) != nil {
		ERROR(w, http.StatusUnauthorized, fmt.Errorf("Wrong password"))
		return
	}

	token, err := auth.CreateToken(dbUser.ID)

	JSON(w, http.StatusOK, token)
}

// validates admin credentials and returns jwt auth_token
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var dbUser models.User
	result := db.DB.Debug().Where("id = ?", 1).First(&dbUser, "email = ?", strings.ToLower(user.Email))
	if result.Error != nil {
		ERROR(w, http.StatusUnauthorized, result.Error)
		return
	}

	if auth.VerifyPassword(dbUser.Password, user.Password) != nil {
		ERROR(w, http.StatusUnauthorized, fmt.Errorf("Wrong password"))
		return
	}

	token, err := auth.CreateToken(dbUser.ID)

	JSON(w, http.StatusOK, token)
}
