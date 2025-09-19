package main

import (
    "fmt"
    "log"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    // Choose a secure password
    password := "123" 
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Password:", password)
    fmt.Println("Hash:    ", string(bytes))
}