package tokenHandler

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type serviceClaims struct {
	Service string
	jwt.StandardClaims
}

type userClaims struct {
	UserID string
	jwt.StandardClaims
}

func HandleVerifyTokenAndGetService(serviceToken string, keyFunc func(*jwt.Token) (interface{}, error)) string {
	s := serviceToken
	c := &serviceClaims{}

	t, err := jwt.ParseWithClaims(s, c, keyFunc)
	if err != nil {
		fmt.Println("Failed because error: ", err)
		return ""
	}
	if !t.Valid {
		fmt.Println("Failed because invalid")
		return ""
	}
	return c.Service
}

func HandleVerifyTokenAndGetUserId(userToken string, keyFunc func(*jwt.Token) (interface{}, error)) map[string]string {
	u := userToken
	c := &userClaims{}

	t, err := jwt.ParseWithClaims(u, c, keyFunc)
	if err != nil {
		fmt.Println("Failed because error:")
		fmt.Println(err)
		return nil
	}
	if !t.Valid {
		fmt.Println("Failed because invalid")
		return nil
	}
	return map[string]string{"userId": c.UserID}
}

func CreateTestServiceToken() string {
	keyfile, _ := ioutil.ReadFile("private.pem")
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyfile)

	c := serviceClaims{
		Service: "front-end",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS512, c)
	s, _ := t.SignedString(key)
	return s
}

func CreateTestUserToken() string {
	keyfile, _ := ioutil.ReadFile("private.pem")
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyfile)

	c := userClaims{
		UserID: "123456",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS512, c)
	s, _ := t.SignedString(key)
	return s
}
