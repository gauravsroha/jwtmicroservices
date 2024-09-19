package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    jwt "github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

// FIXED: Corrected function signature to return (string, error)
func GetJWT() (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    // FIXED: Correctly access claims as MapClaims
    claims := token.Claims.(jwt.MapClaims)
    claims["authorized"] = true
    claims["client"] = "Gaurav"
    claims["aud"] = "billing.jwtgo.io"
    claims["iss"] = "jwtgo.io"
    claims["expiry"] = time.Now().Add(time.Minute*1).Unix()
    
    tokenString, err := token.SignedString(MySigningKey)
    if err != nil {
        // FIXED: Use fmt.Errorf to create an error, don't return it directly
        return "", fmt.Errorf("Something went wrong: %s", err.Error())
    }
    return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
    // FIXED: GetJWT no longer requires arguments
    validToken, err := GetJWT()
    fmt.Println(validToken)
    if err != nil {
        fmt.Println("Failed to generate token")
        // FIXED: Added return to stop execution if there's an error
        return
    }
    fmt.Fprintf(w, validToken) // writes the value of token to w which is the response
}

func handleRequests() {
    http.HandleFunc("/", Index)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    handleRequests()
}