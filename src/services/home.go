package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"database/sql"
	"errors"
)

func (s Service) HomePage(ctx context.Context, user models.User, band models.Band) (models.HomePageData, error) {

	messages, err := database.ChatsTableGetThreeRecentChats(user.UserID, band.BandID)
	if err != nil {
		return models.HomePageData{}, err
	}

	event, err := database.EventsTableGetNextEvent(band.BandID, user.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.HomePageData{}, err
	}

	todoItems, err := database.TodoItemsTableGetNextItemsByUser(user.UserID, band.BandID)
	if err != nil {
		return models.HomePageData{}, err
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(ctx, user.UserID, band.BandID)
	if err != nil {
		return models.HomePageData{}, err
	}

	data := models.HomePageData{
		User:          user,
		Band:          band,
		Event:         event,
		Messages:      messages,
		ToDoItems:     todoItems,
		Notifications: notifications,
	}

	return data, nil
}
