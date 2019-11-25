package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	u "github.com/glorinli/go-jwt-simple-auth/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const TYPE_NORMAL = 0
const TYPE_ADMIN = 100

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

// a struct to rep user account
type Account struct {
	gorm.Model
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Token    string   `json:"token" sql:"-"`
	Type     uint     `json:"type"`
	Profile  *Profile `json:"profile"`
}

// Profile
type Profile struct {
	gorm.Model
	NickName  string `json:"nickname"`
	AvatarUrl string `json:"avatarurl"`
	Age       uint   `json:"age"`
	AccountId uint
}

// validate imcoming user details
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required and at least 6 characters"), false
	}

	temp := &Account{}

	// check for errors and duplicated email
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error, please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email has been taken"), false
	}

	return u.Message(true, "Checking passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	if account.Profile == nil {
		account.Profile = new(Profile)
	}

	err := GetDB().Create(account).Error

	if err != nil || account.ID <= 0 {
		return u.Message(false, "Fail to create account, connection error")
	}

	// Create new JWT token for the new registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = ""
	response := u.MessageWithData(true, "Account has been created", account)
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Where("Email = ?", email).First(account).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error, please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials, please retry")
	}

	account.Password = ""

	// Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	resp := u.MessageWithData(true, "Logged in", account)
	return resp
}

func GetUser(u uint) *Account {
	acc := &Account{}
	err := GetDB().Where("ID = ?", u).First(acc).Error

	if err != nil {
		fmt.Println("Fail to get user: " + err.Error())
		return nil
	}

	fmt.Println("GetUser, email: " + acc.Email)

	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc
}
