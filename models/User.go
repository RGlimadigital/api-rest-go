package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

// Credentials para guaradr datos de autenticación/autorización
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User Struct con los datos del perfil del usuario
type User struct {
	gorm.Model
	UserID int `gorm:"primary_key" json:"user_id"`
	Credentials
	Tasks []Task `gorm:"foreignkey:UserID"` //crear en picture
}

// NewUser Constructor para crear nuevos usuarios
func NewUser(userID int, username string, password string) *User {
	return &User{
		UserID: userID,
		Credentials: Credentials{
			username,
			password,
		},
	}
}

// NewCredentialsJSON Para pasar de JSON del body a objeto User
func NewCredentialsJSON(jsonBytes []byte) *Credentials {
	cred := new(Credentials)
	err := json.Unmarshal(jsonBytes, cred)
	if err == nil {
		return cred
	}
	return nil
}
func NewUserJSON(jsonBytes []byte) *User {
	user := new(User)
	err := json.Unmarshal(jsonBytes, user)
	if err == nil {
		return user
	}
	return nil
}
