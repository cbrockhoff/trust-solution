package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/mettegroenbech/policyGateway/policyTrees"
	"github.com/mettegroenbech/policyGateway/serviceWhitelist"
	"github.com/mettegroenbech/policyGateway/tokenHandler"
)

var policies policyTrees.PolicyTrees

func main() {
	listenPortPtr := flag.String("listenPort", "15151", "the port to listen to")
	forwardAddrPtr := flag.String("forwardAddr", "localhost", "the address to forward the request to")
	forwardPortPtr := flag.String("forwardPort", "14141", "the port to forward the request to")
	policyFilePath := flag.String("policy", "serviceWhitelist.json", "the policy JSON file")
	flag.Parse()

	data, err := ioutil.ReadFile(*policyFilePath)
	if err != nil {
		fmt.Println("Error reading the whitelist file", err)
		return
	}

	var whitelist serviceWhitelist.ServiceWhitelist
	json.Unmarshal([]byte(data), &whitelist)

	policies = policyTrees.BuildPolicyTrees(whitelist)

	fmt.Println("Printing the policy trees:")
	policyTrees.Print(policies)
	fmt.Println()

	keyfile, _ := ioutil.ReadFile("public.pem")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyfile)

	s := &server{
		forwardAddr: *forwardAddrPtr,
		forwardPort: *forwardPortPtr,
		key:         key,
	}
	http.Handle("/", s)
	// log.Fatal(http.ListenAndServe(":"+*listenPortPtr, nil))
	log.Fatal(http.ListenAndServeTLS(":"+*listenPortPtr, "server.crt", "server.key", nil))
}

type server struct {
	forwardAddr string
	forwardPort string
	key         interface{}
}

func (s *server) getKey(_ *jwt.Token) (interface{}, error) {
	return s.key, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serviceToken, foundServiceToken := r.Header["Service-Token"]
	userToken, foundUserToken := r.Header["User-Token"]
	if !foundUserToken {
		userToken = []string{""}
	}
	if foundServiceToken {
		caller := tokenHandler.HandleVerifyTokenAndGetService(serviceToken[0], s.getKey)
		claims := tokenHandler.HandleVerifyTokenAndGetUserId(userToken[0], s.getKey)

		if caller == "" {
			fmt.Println(fmt.Sprintf("Denying %v: %v %v\tNo Service claim found in token", r.RemoteAddr, r.Method, r.URL.Path))
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "request NOT allowed"}`))
			return
		}

		if policyTrees.IsRequestAllowed(policies, r.Method, r.URL.Path, caller, claims) {
			fmt.Println(fmt.Sprintf("Allowing %v: %v %v", caller, r.Method, r.URL.Path))
			s.forwardRequest(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(fmt.Sprintf("Denying %v: %v %v", caller, r.Method, r.URL.Path))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "request NOT allowed"}`))
		}
	} else {
		fmt.Println(fmt.Sprintf("Denying %v: %v %v \tNo service token", r.RemoteAddr, r.Method, r.URL.Path))
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "request NOT allowed"}`))
	}
}

func (s *server) forwardRequest(w http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse("http://" + s.forwardAddr + ":" + s.forwardPort)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = s.forwardAddr + ":" + s.forwardPort
	req.URL.Scheme = "http"
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = s.forwardAddr + ":" + s.forwardPort

	fmt.Println(fmt.Sprintf("Proxying %v %v", req.Method, req.URL.Path))
	proxy.ServeHTTP(w, req)
}
