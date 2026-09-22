package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (s Service) RegistrationCreateUserProfile(firstName, lastName, displayName, timezone, accessCode string) (models.UserRegistration, error) {
	log.Println("- RegistrationCreateUserProfile")

	newID := uuid.New().String()

	var user models.UserRegistration

	user.UserRegistrationID = newID

	newUser, err := database.UsersRegTableCreateInititialUser(user)
	if err != nil {
		log.Println("   Unable to create initial user registration")
		return models.UserRegistration{}, err
	}

	return newUser, nil
}

func (s Service) RegistrationSavePassword(registrationID, password string) (models.UserRegistration, error) {
	log.Println("- RegistrationSavePassword")

	user, err := database.UsersRegTableGetUserByID(registrationID)
	if err != nil {
		log.Println("Unable to get registration user by ID")
		return models.UserRegistration{}, err
	}

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return models.UserRegistration{}, err
	}

	fmt.Println("passwordHash: ", passwordHash)

	user.PasswordHash = passwordHash

	return user, nil
}
