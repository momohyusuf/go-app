package authservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/momoh-yusuf/note-app/config"
	"github.com/momoh-yusuf/note-app/generated_sql"
	"github.com/momoh-yusuf/note-app/utils"
	"golang.org/x/crypto/bcrypt"
)

func HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var RequestBody utils.RegisterUserBody
	// extract the json field in the r.body with decoder
	decoder := json.NewDecoder(r.Body)
	// check if an error occurs while parsing
	err := decoder.Decode(&RequestBody)

	data := map[string]interface{}{
		"msg": "",
	}

	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !RequestBody.ValidateUserName() {
		data["msg"] = "Please Provide a valid user name"
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	if !utils.ValidateEmail(RequestBody.Email) {
		data["msg"] = fmt.Sprintf("%v is not valid email address. Provide a valid email", RequestBody.Email)
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	if !utils.ValidateStrongPassword(RequestBody.Password) {
		data["msg"] = fmt.Sprintf("%v is not strong password. Provide a strong password", RequestBody.Password)
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	if RequestBody.Age < 10 {
		data["msg"] = "Sorry you're too young to use this app. Must be older than 9 years"
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	// insert into the data base

	result, err := config.Db_Query().CreateUser(ctx, generated_sql.CreateUserParams{
		UserName:     strings.TrimSpace(RequestBody.User_name),
		Email:        RequestBody.Email,
		UserPassword: RequestBody.HashUserPassword(),
		Age:          RequestBody.Age,
		CreatedAt:    pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		utils.CustomResponseInJson(w, http.StatusInternalServerError, data, err)
		return
	}
	utils.CustomResponseInJson(w, http.StatusCreated, result, nil)
	w.WriteHeader(201)
	w.Write([]byte("Registered successfully"))
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleUerLogin(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background()
	var requestBody LoginRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)

	data := map[string]interface{}{}
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if !utils.ValidateEmail(requestBody.Email) {
		data["msg"] = fmt.Sprintf("%v is not valid email address. Provide a valid email", requestBody.Email)
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	if !utils.ValidateStrongPassword(requestBody.Password) {
		data["msg"] = fmt.Sprintf("%v is not strong password. Provide a strong password", requestBody.Password)
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	// validate that both email and password are valid
	user, err := config.Db_Query().FindUserByEmail(context.Background(), requestBody.Email)

	if err != nil {
		data["msg"] = fmt.Sprintf("No record found for %v", requestBody.Email)
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	// Compare user password with the hashPassword

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(requestBody.Password))
	if err != nil {
		data["msg"] = "Email or password incorrect"
		utils.CustomResponseInJson(w, http.StatusBadRequest, data, nil)
		return
	}

	data = jwt.MapClaims{
		"user_email": user.Email,
		"user_id":    user.UserID,
		"user_type":  user.UserType,
		"exp":        time.Now().Add(10 * time.Hour),
	}
	accessToken := utils.GenerateJwtToken(data)

	data["token"] = accessToken
	data["msg"] = "Success"
	utils.CustomResponseInJson(w, http.StatusOK, data, nil)
}
