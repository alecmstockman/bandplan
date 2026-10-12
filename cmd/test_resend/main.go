package main

import (
	"bandplan/src/clients"
	"log"
)

func main() {
	log.Println("Testing Resend API")

	if err := clients.SendTestVerificationEmail(); err != nil {
		log.Printf("unable to send test email: %v", err)
	}
}
