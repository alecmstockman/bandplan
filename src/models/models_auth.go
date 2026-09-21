package models

import "time"

type userRegistration struct {
	ID     string
	BandID string

	BandName string

	FirstName    string
	LastName     string
	Email        string
	PasswordHash string

	EmailValidated bool
	ValidationSent time.Time
	validated      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
