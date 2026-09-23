package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (s Service) RegistrationSaveUserProfile(registrationID, email, firstName, lastName, displayName, timezone, accessCode string) (models.UserRegistration, error) {
	log.Println("- RegistrationSaveUserProfile")

	user := models.UserRegistration{
		UserRegistrationID: registrationID,
	}
	user.FirstName = firstName
	user.LastName = lastName
	user.DisplayName = displayName
	user.Email = email
	user.Timezone = timezone
	user.AccessCodeHash = accessCode

	if registrationID != "" {
		updatedUser, err := database.UsersRegTableUpdateInitialUser(user)
		if err != nil {
			log.Println("   Unable to update initial user registration")
			return models.UserRegistration{}, err
		}
		return updatedUser, nil
	}

	user.UserRegistrationID = uuid.New().String()

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
