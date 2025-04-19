package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte (os.Getenv("SECRET_KEY"))

func CreateToken (userId string) (string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"userId":userId,
		"exp":time.Now().Add(time.Hour*24).Unix(),
	})
	tokenString ,err:= token.SignedString(secretKey)
	if err != nil{
		return "",err
	}
	return tokenString,nil
}
func VerifyToken(tokenString string) (string , error){
	token, err := jwt.Parse(tokenString ,func (token *jwt.Token) (any , error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,errors.New("unexpected signing method")
		}
		return secretKey,nil
	})
	if err!=nil || !token.Valid {
		return "",errors.New("invalid token")
	}
	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("could not parse claims")

	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userId, nil
}

