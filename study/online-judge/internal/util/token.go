package util

import (
	"github.com/dgrijalva/jwt-go"
)

type UerClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var key = []byte("gin-gorm-oj-key")

// GenerateToken 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UerClaims := &UerClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UerClaims)
	str, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return str, nil
}

// AnalyzeToken 解析token
func AnalyzeToken(str string) (*UerClaims, error) {
	userClaim := new(UerClaims)
	claims, err := jwt.ParseWithClaims(str, userClaim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, err
	}
	return userClaim, nil
}
