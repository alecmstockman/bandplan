package models

import "time"

type LoginPageData struct {
	CSRFToken string
}

type RegistrationPages struct {
	CSRFToken         string
	User              UserRegistration
	Band              Band
	AccessCode        string
	RegistrationToken string
	BandName          string
}

type UserRegistration struct {
	CSRFToken             string
	ID                    string
	RegistrationToken     string
	RegistrationTokenHash string
	AccessCodeHash        string

	FirstName   string
	LastName    string
	DisplayName string

	BandID        string
	Timezone      string
	Email         string
	EmailVerified bool
	PasswordHash  string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
