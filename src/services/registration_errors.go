package services

import (
	"bandplan/src/database"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type RegistrationErrorKind uint8

const (
	RegistrationInvalid RegistrationErrorKind = iota + 1
	RegistrationForbidden
	RegistrationNotFound
	RegistrationConflict
	RegistrationExpired
	RegistrationRateLimited
	RegistrationInternal
)

type RegistrationError struct {
	Kind RegistrationErrorKind
	Err  error
}

func (e *RegistrationError) Error() string {
	return e.Err.Error()
}

func (e *RegistrationError) Unwrap() error {
	return e.Err
}

func registrationError(kind RegistrationErrorKind, err error) error {
	return &RegistrationError{Kind: kind, Err: err}
}

func classifyRegistrationError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, database.ErrRegistrationExpired) || errors.Is(err, database.ErrAccessCodeExpired) {
		return registrationError(RegistrationExpired, err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return registrationError(RegistrationNotFound, err)
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return registrationError(RegistrationConflict, err)
	}

	return registrationError(RegistrationInternal, err)
}
