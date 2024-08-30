package main

import (
	"fmt"
	"log"

	"github.com/roborobs1023/tools/validate"
)

type User struct {
	Name   string `validate:"min=2,max64"`
	Email  string `validate:"email,nonDisposable"`
	Domain string `validate:"required,domain,nonDisposable"`
}

func main() {
	users := []User{
		{
			Name:   "A",
			Email:  "testing@go.dev",
			Domain: "testing.com",
		},
		{
			Name:   "Austin",
			Email:  "testing@example.com",
			Domain: "apple",
		},
		{
			Name:   "Apple",
			Email:  "apple@gmail.com",
			Domain: "google.com",
		},
		{
			Name:  "Apple",
			Email: "apple@gmail.com",
		},
	}

	for i, u := range users {
		fmt.Print()
		err := validate.Validate(u, &validate.Config{})

		if err != nil {
			fmt.Println("User", i+1)
			log.Println(err)
		} else {
			fmt.Println("User", i+1, "Valid")
		}

	}
}
