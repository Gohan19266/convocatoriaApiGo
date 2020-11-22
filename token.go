package main

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/jwt.pem"
	pubKeyPath  = "keys/jwt.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type accessDetails struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PhotoURL string `json:"photoUrl"`
	Third    bool   `json:"third"`
	Code     string `json:"code,omitempty"`
	Created  string `json:"created"`
}

type errorResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				w.WriteHeader(http.StatusUnauthorized)
				response := errorResponse{false, "Token expirado!"}
				json.NewEncoder(w).Encode(response)
				return err
			}
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			response := errorResponse{false, "Token invalido!"}
			json.NewEncoder(w).Encode(response)
			return err
		}
	}
	return nil
}

func extractTokenMetadata(w http.ResponseWriter, r *http.Request) (*accessDetails, error) {
	e := verifyToken(w, r)
	if e != nil {
		return nil, e
	}
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		info, ok := claims["user_claims"]
		if !ok {
			return nil, err
		}
		jsonString, _ := json.Marshal(info)
		s := &accessDetails{}
		json.Unmarshal(jsonString, s)
		return s, nil
	}
	return nil, err
}
