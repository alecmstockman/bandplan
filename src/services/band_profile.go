package services

import (
	"bandplan/src/clients"
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
	"time"
)

func (s Service) BandSendEmailVerification(ctx context.Context, band models.Band) error {
	if band.EmailAddress == "" {
		return nil
	}

	token, err := helpers.GenerateSessionToken()
	if err != nil {
		return fmt.Errorf("generate band email verification token: %w", err)
	}
	tokenHash := helpers.HashRegistrationCode(token)

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin band email verification transaction: %w", err)
	}
	defer tx.Rollback()

	if err := database.BandEmailVerificationTableCreateCode(
		ctx,
		tx,
		band.BandID,
		band.EmailAddress,
		tokenHash,
		time.Now().UTC().Add(30*time.Minute),
	); err != nil {
		return err
	}

	if err := clients.SendBandEmailVerification(models.BandVerificationPage{
		BandID:   band.BandID,
		BandName: band.Name,
		Email:    band.EmailAddress,
		Token:    token,
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit band email verification transaction: %w", err)
	}
	return nil
}

func (s Service) BandResendEmailVerification(ctx context.Context, band models.Band) error {
	if band.EmailAddress == "" || band.IsEmailVerified {
		return nil
	}

	err := s.BandSendEmailVerification(ctx, band)
	if errors.Is(err, database.ErrEmailVerificationRateLimited) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("resend band email verification: %w", err)
	}

	return nil
}
