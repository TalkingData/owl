package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func initPublicKeyAndPrivateKey() (err error) {
	signBytes, err := ioutil.ReadFile(config.PrivateKeyFile)
	if err != nil {
		return
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return
	}

	verifyBytes, err := ioutil.ReadFile(config.PublicKeyFile)
	if err != nil {
		return
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return
	}
	return
}

func GenerateToken(email string, user *User) (string, error) {
	// Create the Claims
	if user == nil {
		user = &User{}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      user.Username,
		"username": user.DisplayName,
		"phone":    user.PhoneNum,
		"email":    email,
		"exp":      time.Now().Add(time.Duration(config.TokenExpired) * time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})
	return token.SignedString(signKey)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt token error %s", err.Error())
	}

	return token, nil
}
