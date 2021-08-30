package models

import (
	"crypto/rand"
	"database/sql"
	"depmod/db"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID       int64
	Username string
	Email    string
	Password string
	Token    sql.NullString
}

// func FetchAllUser() Response {

// 	var users []Users
// 	var res Response

// 	con := db.CreateCon()

// 	con.Find(&users)

// 	res.Status = http.StatusOK
// 	res.Message = "success"
// 	res.Data = users

// 	return res

// }

func RegisterUser(username string, email string, password string) Response {
	var res Response
	con := db.CreateCon()

	hash, _ := HashPassword(password)

	user := Users{
		Username: username,
		Email:    email,
		Password: password,
		Token:    sql.NullString{Valid: false}}
	res = UserRegisterValidator(&user, res, con)
	user.Password = hash
	res.Data = user

	if res.Status == 200 {
		con.Create(&user)
	}
	return res

}

func LoginUser(email string, password string) Response {
	var res Response
	con := db.CreateCon()
	user, found := UserLoginValidator(con, email, password)
	if found {
		res.Status = http.StatusOK
		res.Message = "Success"
		user.Token = sql.NullString{String: tokenGenerator(), Valid: true}
		res.Data = user
		con.Save(&user)
		return res
	}

	res.Status = http.StatusBadRequest
	res.Message = "your Password/Email is wrong"

	return res

}

func LogoutUser(token string) Response {
	var res Response
	var user Users
	a := strings.Split(token, " ")[1]
	con := db.CreateCon()
	if con.Where("token = ?", a).First(&user).RowsAffected == 0 {
		res.Status = http.StatusBadRequest
		res.Message = "Session not found"

		return res
	}
	user.Token = sql.NullString{Valid: false}
	con.Save(&user)
	res.Status = http.StatusOK
	res.Message = "Success"

	return res
}

func UserLoginValidator(con *gorm.DB, email string, password string) (Users, bool) {
	var user Users
	if con.Where("email = ?", email).First(&user).RowsAffected > 0 && CheckPasswordHash(password, user.Password) {
		return user, true
	}
	return user, false

}

func tokenGenerator() string {
	b := make([]byte, 64)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func UserRegisterValidator(user *Users, res Response, con *gorm.DB) Response {

	if con.Where("email = ?", user.Email).First(&user).RowsAffected > 0 {
		res.Status = http.StatusBadRequest
		res.Message = "email already registered"
		return res
	}

	err := validation.Validate(
		user.Username,
		validation.Required,
		validation.Length(4, 32))

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "username: " + err.Error()
		return res
	}

	err = validation.Validate(
		user.Email,
		validation.Required,
		validation.Length(4, 64),
		is.Email)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "email: " + err.Error()
		return res
	}
	err = validation.Validate(
		user.Password,
		validation.Required,
		validation.Length(6, 32))

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "password: " + err.Error()
		return res
	}
	res.Status = http.StatusOK
	res.Message = "Success"
	return res
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
