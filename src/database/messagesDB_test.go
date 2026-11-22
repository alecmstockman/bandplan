package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMessagesTableCreateMessageReturnsReplyPreview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create database mock: %v", err)
	}

	originalDB := DB
	DB = db
	t.Cleanup(func() {
		DB = originalDB
		db.Close()
	})

	createdAt := time.Now()
	mock.ExpectQuery(`WITH reply_target AS`).
		WithArgs(
			"band-id",
			"user-id",
			"chat-id",
			"message-id",
			"Reply body",
			"parent-message-id",
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"message_id",
			"band_id",
			"user_id",
			"chat_id",
			"body",
			"reply_to_message_id",
			"reply_to_user_name",
			"reply_to_body",
			"created_at",
		}).AddRow(
			1,
			"message-id",
			"band-id",
			"user-id",
			"chat-id",
			"Reply body",
			"parent-message-id",
			"Parent User",
			"Original body",
			createdAt,
		))

	message, err := MessagesTableCreateMessage(
		"band-id",
		"user-id",
		"Current User",
		"chat-id",
		"message-id",
		"Reply body",
		"parent-message-id",
		"America/New_York",
	)
	if err != nil {
		t.Fatalf("create reply: %v", err)
	}
	if message.ReplyToMessageID == nil || *message.ReplyToMessageID != "parent-message-id" {
		t.Fatalf("reply message ID = %v; want parent-message-id", message.ReplyToMessageID)
	}
	if message.ReplyToUserName != "Parent User" {
		t.Fatalf("reply user name = %q; want Parent User", message.ReplyToUserName)
	}
	if message.ReplyToBody != "Original body" {
		t.Fatalf("reply body = %q; want Original body", message.ReplyToBody)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
