package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"context"
	"log"
	"log/slog"
	"unicode/utf8"

	"github.com/google/uuid"
)

func (s Service) RegistrationLoadAccessData(ctx context.Context, registrationID, accessCode string) (models.RegistrationPages, error) {

	if accessCode != "" {
		_, err := database.AccessCodesTableValidateCodeReturnBandID(ctx, accessCode)
		if err != nil {
			return models.RegistrationPages{}, err
		}
	}

	var user models.UserRegistration

	if registrationID != "" {
		registeredUser, err := database.UsersRegTableGetUserRegistrationID(registrationID)
		if err != nil {
			return models.RegistrationPages{}, err
		}
		user = registeredUser
	}

	data := models.RegistrationPages{
		AccessCode:     accessCode,
		RegistrationID: registrationID,
		User:           user,
	}

	return data, nil
}

func (s Service) RegistrationSaveUserProfile(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("- RegistrationSaveUserProfile")

	if user.UserRegistrationID != "" {
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

	user, err := database.UsersRegTableGetUserRegistrationID(registrationID)
	if err != nil {
		log.Println("Unable to get registration user by ID")
		return models.UserRegistration{}, err
	}

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return models.UserRegistration{}, err
	}

	user.PasswordHash = passwordHash

	return user, nil
}

func (s Service) RegistrationBandPageSubmit(ctx context.Context, registrationID, bandName string) (models.RegistrationPages, error) {
	if _, err := uuid.Parse(registrationID); err != nil {
		return models.RegistrationPages{}, err
	}
	valid, err := database.UserRegTableValidateRegistrationID(registrationID)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(ctx),
			"error", err,
		)
		return models.RegistrationPages{}, err
	}

	if valid == false {
		return models.RegistrationPages{}, err
	}

	if bandName == "" || utf8.RuneCountInString(bandName) > 100 {
		return models.RegistrationPages{}, err
	}

	data := models.RegistrationPages{
		RegistrationID: registrationID,
		BandName:       bandName,
	}

	return data, nil
}
