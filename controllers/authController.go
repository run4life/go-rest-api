package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/glorinli/go-jwt-simple-auth/models"
	"github.com/glorinli/go-jwt-simple-auth/utils"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}

	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid info: "+err.Error()))
		return
	}

	utils.Respond(w, account.Create())
}

var Login = func(w http.ResponseWriter, r *http.Request) {
	email := utils.GetRequestParam(r, "email")

	pwd := utils.GetRequestParam(r, "password")

	if email == "" || pwd == "" {
		utils.Respond(w, utils.Message(false, "Invalid request param, email: "+email))
		return
	}

	utils.Respond(w, models.Login(email, pwd))
}

var Me = func(w http.ResponseWriter, r *http.Request) {
	value := r.Context().Value("user")

	fmt.Println("Context.user", value)

	userId, ok := value.(uint)

	if !ok {
		utils.Respond(w, utils.Message(false, "Invalid userId: "+string(userId)))
		return
	}

	account := models.GetUser(userId)

	if account == nil {
		utils.Respond(w, utils.Message(false, "User not found"))
		return
	}

	message := utils.MessageWithData(true, "Success", account)

	utils.Respond(w, message)
}
