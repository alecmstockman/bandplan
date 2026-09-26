package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (s Service) RegistrationSaveUserProfile(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("- RegistrationSaveUserProfile")

	if user.UserRegistrationID != "" {
		fmt.Println("RegistrationID present: ", user.UserRegistrationID)
		updatedUser, err := database.UsersRegTableUpdateInitialUser(user)
		if err != nil {
			log.Println("   Unable to update initial user registration")
			return models.UserRegistration{}, err
		}
		return updatedUser, nil
	}

	user.UserRegistrationID = uuid.New().String()

	fmt.Printf("\nUser: %+v\n", user)

	newUser, err := database.UsersRegTableCreateInititialUser(user)
	fmt.Println("creating registration ID ")
	if err != nil {
		log.Println("   Unable to create initial user registration")
		return models.UserRegistration{}, err
	}

	return newUser, nil
}

func (s Service) RegistrationSavePassword(registrationID, password string) (models.UserRegistration, error) {
	log.Println("- RegistrationSavePassword")

	fmt.Println("registrationID: ", registrationID)

	user, err := database.UsersRegTableGetUserByID(registrationID)
	if err != nil {
		log.Println("Unable to get registration user by ID")
		return models.UserRegistration{}, err
	}

	fmt.Println("User: ", user)

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return models.UserRegistration{}, err
	}

	fmt.Println("passwordHash: ", passwordHash)

	user.PasswordHash = passwordHash

	return user, nil
}
