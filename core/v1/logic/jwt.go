package logic

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"log"
)

type jwtService struct {
	jwt v1.Jwt
}

func (j jwtService) ValidateToken(tokenString string) (bool, *jwt.Token) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwt.PublicKey, nil
	})
	if err != nil {
		log.Print("[ERROR]: Token is invalid! ", err.Error())
		return false, nil
	}
	return true, token

}

func getPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(config.Publickey))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		log.Print(err.Error())
		panic(err)
	}
	return publicKeyImported
}

// NewJwtService returns Jwt type service
func NewJwtService() service.Jwt {
	return jwtService{
		jwt: v1.Jwt{
			PublicKey: getPublicKey(),
		},
	}
}
