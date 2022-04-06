package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UerClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var key = []byte("gin-gorm-oj-key")

func TestGenerateToken(t *testing.T) {
	UerClaims := &UerClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UerClaims)
	str, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE
	fmt.Println(str)
}

func TestAnalyzeToken(t *testing.T) {
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE"
	userClaim := new(UerClaims)
	claims, err := jwt.ParseWithClaims(str, userClaim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}
