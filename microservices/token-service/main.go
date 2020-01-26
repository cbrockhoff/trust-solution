package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

type tokenResponse struct {
	Token string
}

type claims struct {
	Service string
	jwt.StandardClaims
}

func handleGetToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params, keyFilePath string) {
	authHeaders := r.Header["Authorization"]
	if len(authHeaders) == 0 {
		log.Println(fmt.Sprintf("[Client %v]: Authorization header missing, returning 401", r.RemoteAddr))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	headerParts := strings.Split(authHeaders[0], " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "basic" {
		log.Println(fmt.Sprintf("[Client %v]: Authorization header value %v malformed, returning 401", r.RemoteAddr, headerParts))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(headerParts[1])
	if err != nil {
		log.Println(fmt.Sprintf("[Client %v]: Authorization header value %v could not be decoded - see error below, returning 401", r.RemoteAddr, headerParts[1]))
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	authParts := strings.Split(string(decoded), ":")
	username, password := authParts[0], authParts[1]

	lookupTime := rand.Intn(200) + 50
	time.Sleep(time.Duration(lookupTime) * time.Millisecond) // Ooooh, doing database lookups!

	if password != username+"_password" {
		log.Println(fmt.Sprintf("[Client %v]: Wrong password, returning 401", r.RemoteAddr))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	keyfile, _ := ioutil.ReadFile(keyFilePath)
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyfile)

	c := claims{
		Service: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, c)
	s, _ := t.SignedString(key)
	log.Println(fmt.Sprintf("[Client %v]: Created token for user %v", r.RemoteAddr, username))
	json.NewEncoder(w).Encode(tokenResponse{Token: s})
}

func main() {
	fmt.Println("Starting service...")
	listenPort := flag.String("port", "80", "The port to listen on")
	privateKeyFile := flag.String("key-path", "private.pem", "Pem file with private key to sign tokens with")
	flag.Parse()

	r := httprouter.New()
	r.GET("/tokens", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handleGetToken(w, r, ps, *privateKeyFile)
	})
	log.Fatal(http.ListenAndServe(":"+*listenPort, r))
	fmt.Println("Quitting service...")
}
