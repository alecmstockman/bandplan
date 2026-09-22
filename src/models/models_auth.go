package models

import "time"

type RegistrationPages struct {
	User UserRegistration
}

type UserRegistration struct {
	ID                 string
	UserRegistrationID string
	AccessCodeHash     string

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
