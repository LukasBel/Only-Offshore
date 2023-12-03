package Handlers

import (
	"github.com/joho/godotenv"
	"net/smtp"
	"os"
)

type EmailAgent struct {
	from     string
	password string
}

func SendMail(to []string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	agent := EmailAgent{from: os.Getenv("FROM"), password: os.Getenv("PASSWORD")}
	auth := smtp.PlainAuth("", agent.from, agent.password, "smtp.gmail.com")

	err = smtp.SendMail("smtp.gmail.com:587", auth, agent.from, to, Message())
	if err != nil {
		return err
	}

	return nil
}

func Message() []byte {
	subject := "Gym Stats Update!"
	body := ""
	message := []byte(subject + body)
	return message
}
