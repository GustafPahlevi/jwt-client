package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(`mysuperphrase`)

var URLServer = "http://localhost:8080/"

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "gustaf"
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	fmt.Println("token is: ", tokenString)
	return tokenString, nil
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := GenerateJWT()
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Fprintf(w, HTTPCall(token))
}

func main() {
	fmt.Println("starting simple client")
	http.HandleFunc("/", handleHTTP)
	_ = http.ListenAndServe(":8081", nil)
}

func HTTPCall(token string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, URLServer, nil)
	if err != nil {
		fmt.Println("err is: ", err)
		return ""
	}
	req.Header.Set("Token", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		return ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err is: ", err)
		return ""
	}

	fmt.Println("response is: ", string(body))
	return string(body)
}
