package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword uses bcrypt to generate a hash from a password.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash compares a password with its hash.
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}