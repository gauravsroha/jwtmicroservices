package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Your jwt token is secured and verified")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method")
				}
				// FIXED: Correctly type assert token.Claims to jwt.MapClaims
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					return nil, fmt.Errorf("invalid claims")
				}
				aud := "billing.jwtgo.io"
				// FIXED: Use claims.VerifyAudience instead of token.Claims.(jwt.MapClaims).VerifyAudience
				checkAudience := claims.VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf("invalid audience")
				}
				iss := "jwtgo.io"
				// FIXED: Use claims.VerifyIssuer instead of token.Claims.(jwt.MapClaims).VerifyISS
				checkISS := claims.VerifyIssuer(iss, false)
				if !checkISS {
					return nil, fmt.Errorf("invalid issuer")
				}
				return MySigningKey, nil
			})
			if err != nil {
				fmt.Fprint(w, err.Error())
				// FIXED: Added return statement to prevent further execution if there's an error
				return
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprint(w, "No authorization token provided")
		}
	})
}

func handleRequests() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	// FIXED: Changed Printf to Println for better output, and improved the message
	fmt.Println("Server starting...")
	handleRequests()
}
