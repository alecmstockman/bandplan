package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"unicode/utf8"
)

func (s Service) RegistrationLoadAccessData(ctx context.Context, registrationToken, accessCode string) (models.RegistrationPages, error) {

	if accessCode != "" {

		accessCodeHash := helpers.HashRegistrationCode(accessCode)

		_, err := database.AccessCodesTableValidateCodeReturnBandID(ctx, accessCodeHash)
		if err != nil {
			return models.RegistrationPages{}, classifyRegistrationError(err)
		}
	}

	var user models.UserRegistration

	if registrationToken != "" {

		valid := helpers.ValidateTokenLength(registrationToken)
		if valid == false {
			return models.RegistrationPages{}, registrationError(RegistrationInvalid, errors.New("invalid registration id"))
		}

		registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

		registeredUser, err := database.UsersRegTableGetRegistrationToken(registrationTokenHash)
		if err != nil {
			return models.RegistrationPages{}, classifyRegistrationError(err)
		}
		user = registeredUser

	}

	data := models.RegistrationPages{
		AccessCode:        accessCode,
		RegistrationToken: registrationToken,
		User:              user,
	}

	return data, nil
}

func (s Service) RegistrationSaveUserProfile(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("- RegistrationSaveUserProfile")

	if user.RegistrationToken != "" {
		valid := helpers.ValidateTokenLength(user.RegistrationToken)
		if valid != true {
			return models.UserRegistration{}, registrationError(RegistrationInvalid, errors.New("invalid registration token"))
		}

		registrationToken := user.RegistrationToken
		user.RegistrationTokenHash = helpers.HashRegistrationCode(registrationToken)

		existingUser, err := database.UsersRegTableGetRegistrationToken(user.RegistrationTokenHash)
		if err != nil {
			return models.UserRegistration{}, classifyRegistrationError(err)
		}
		user.AccessCodeHash = existingUser.AccessCodeHash
		user.BandID = existingUser.BandID

		updatedUser, err := database.UsersRegTableUpdateInitialUser(user)
		if err != nil {
			log.Println("   Unable to update initial user registration")
			return models.UserRegistration{}, classifyRegistrationError(err)
		}
		updatedUser.RegistrationToken = registrationToken
		return updatedUser, nil
	}

	newRegistrationToken, err := helpers.GenerateSessionToken()
	if err != nil {
		return models.UserRegistration{}, registrationError(RegistrationInternal, err)
	}

	newRegistrationTokenHash := helpers.HashRegistrationCode(newRegistrationToken)

	user.RegistrationTokenHash = newRegistrationTokenHash

	err = database.UsersRegTableDeleteExpiredUserByEmail(user.Email)
	if err != nil {
		return models.UserRegistration{}, classifyRegistrationError(err)
	}

	newUser, err := database.UsersRegTableCreateInititialUser(user)
	if err != nil {
		log.Println("   Unable to create initial user registration")
		return models.UserRegistration{}, classifyRegistrationError(err)
	}

	newUser.RegistrationToken = newRegistrationToken

	return newUser, nil
}

func (s Service) RegistrationSavePassword(registrationToken, password string) (models.UserRegistration, error) {
	log.Println("- RegistrationSavePassword")

	valid := helpers.ValidateTokenLength(registrationToken)
	if valid == false {
		return models.UserRegistration{}, registrationError(RegistrationInvalid, errors.New("invalid registration id"))
	}

	registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

	user, err := database.UsersRegTableGetRegistrationToken(registrationTokenHash)
	if err != nil {
		return models.UserRegistration{}, classifyRegistrationError(err)
	}

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return models.UserRegistration{}, registrationError(RegistrationInternal, err)
	}

	user.PasswordHash = passwordHash

	return user, nil
}

func (s Service) RegistrationBandPageSubmit(ctx context.Context, registrationToken, bandName string) (models.RegistrationPages, error) {

	valid := helpers.ValidateTokenLength(registrationToken)
	if valid == false {
		return models.RegistrationPages{}, registrationError(RegistrationInvalid, errors.New("invalid registration id"))
	}

	registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

	valid, err := database.UserRegTableValidateRegistrationToken(registrationTokenHash)
	if err != nil {
		return models.RegistrationPages{}, classifyRegistrationError(err)
	}

	if valid == false {
		return models.RegistrationPages{}, registrationError(RegistrationNotFound, errors.New("registration not found"))
	}

	if bandName == "" || utf8.RuneCountInString(bandName) > 100 {
		return models.RegistrationPages{}, registrationError(RegistrationInvalid, errors.New("invalid band name entry"))
	}

	data := models.RegistrationPages{
		RegistrationToken: registrationToken,
		BandName:          bandName,
	}

	return data, nil
}

func (s Service) RegistrationCreateAccessCode(ctx context.Context, user models.User, band models.Band) (string, error) {
	if user.IsAdmin != true {
		return "", registrationError(RegistrationForbidden, errors.New("user is not admin"))
	}

	code := helpers.GenerateAccessCode()
	codeHash := helpers.HashRegistrationCode(code)

	expiresAt := time.Now().Add(24 * time.Hour).UTC()

	err := database.AccessCodesTablesCreateCode(band.BandID, user.UserID, codeHash, expiresAt)
	if err != nil {
		return "", classifyRegistrationError(err)
	}

	html := fmt.Sprintf(`
		<div class="admin-access-code-box">
            <span class="admin-access-code">%v</span>
			<button
              type="button"
              class="admin-display-field-copy"
              onclick="copyToClipboard('%s', this)">
			  <span class="copy-icon">
				  <svg
					xmlns="http://www.w3.org/2000/svg"
					width="20" height="20"
					viewBox="0 0 24 24" fill="none"
					stroke="currentColor" stroke-width="2"
					stroke-linecap="round" stroke-linejoin="round"
					class="lucide lucide-copy-icon lucide-copy">
					<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/>
					<path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
				</svg>
              </span>

			  <span class="check-icon">
				  <svg
					xmlns="http://www.w3.org/2000/svg"
					width="20" height="20"
					viewBox="0 0 24 24" fill="none"
					stroke="currentColor" stroke-width="2"
					stroke-linecap="round" stroke-linejoin="round"
					class="lucide lucide-check-icon lucide-check">
					<path d="M20 6 9 17l-5-5"/>
				</svg>
              </span>
            </button>
		</div>
	`, code, code)

	return html, nil
}
