package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$yGK/Eq40Dj0dO9krcbkvceB9cZOlGKSREQEAeFobeNPX96mvQv4G2"
	password := "admin123"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Comparison failed: %v\n", err)
	} else {
		fmt.Println("Comparison succeeded!")
	}
}
