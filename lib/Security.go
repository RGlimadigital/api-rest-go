package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/RGlimadigital/Tareas-Go/data"
	"github.com/RGlimadigital/Tareas-Go/models"
	"github.com/jinzhu/gorm"
)

type JSONToken struct {
	Token string `json:"token"`
}

type TokenJWT struct {
	UserID int
	jwt.StandardClaims
}

// ValidateCredent Buscar si existe un usuario para estos credenciales
func ValidateCredent(cred *models.Credentials, db *gorm.DB) *models.User {
	user := &models.User{}
	db.Where("username = ? AND password = ?", cred.Username, cred.Password).Find(user)
	if user.UserID > 0 {
		return user
	} else {
		return nil
	}
}

func CreateJWT(usr *models.User) (string, error) {
	loginToken := &TokenJWT{UserID: usr.UserID}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), loginToken) //Libreria
	return jwtToken.SignedString([]byte(os.Getenv("secret_key")))

}

func CreateToken(usr *models.User, cacheClient data.CacheProvider) (token string, err error) {
	h := hmac.New(sha256.New, []byte(os.Getenv("secret_key")))
	_, err = h.Write([]byte(usr.Username))
	if err == nil {
		token = hex.EncodeToString(h.Sum(nil))
		cacheClient.SetExpiration(token, usr, time.Hour)
	}
	return
}

func GetUserJWT(tokenString string, db *gorm.DB) *models.User {
	tokenStruct := new(TokenJWT)
	token, err := jwt.ParseWithClaims(tokenString, tokenStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret_key")), nil
	})
	if err == nil && token.Valid && tokenStruct.UserID > 0 {
		user := new(models.User)
		db.First(user, tokenStruct.UserID)
		return user
	}
	return nil

}

func GetUserTokenCache(tokenString string, cacheClient data.CacheProvider) *models.User {
	if validUser, exists := cacheClient.Get(tokenString); exists && validUser != nil {
		return validUser.(*models.User)
	}
	return nil
}
