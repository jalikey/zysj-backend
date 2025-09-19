package auth

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// A secret key for signing the tokens. In production, use a secure key from env variables.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT creates a new JWT for a given username.
func GenerateJWT(username string) (string, error) {
    if len(jwtSecret) == 0 {
        return "", fmt.Errorf("JWT_SECRET environment variable not set")
    }

    // Create a new token object, specifying signing method and the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
        "iat":      time.Now().Unix(),                      // Issued at
    })

    // Sign and get the complete encoded token as a string using the secret
    tokenString, err := token.SignedString(jwtSecret)
    return tokenString, err
}

// ValidateJWT parses and validates a token string.
func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })
}