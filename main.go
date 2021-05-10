package main

import (
	"errors"
	"github.com/msteinert/pam"
	"log"
)

type Credentials struct {
	User     string
	Password string
}

func (c Credentials) RespondPAM(s pam.Style, msg string) (string, error) {
	switch s {
	case pam.PromptEchoOn:
		return c.User, nil
	case pam.PromptEchoOff:
		return c.Password, nil
	}
	return "", errors.New("unexpected")
}

func main() {
	println("Launching Authentication Server...")

	c := Credentials{
		User:     "louis",
		Password: "",
	}
	transaction, err := pam.Start("", "", c)
	if err != nil {
		log.Fatal(err)
	}
	err = transaction.Authenticate(0)
	if err != nil {
		log.Fatal(err)
	}
}
