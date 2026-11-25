package helpers

import (
	"bandplan/src/models"
)

func ReturnRecurrence(EventRecurrence string) models.EventRecurrence {

	switch EventRecurrence {
	case "daily":
		return models.EventRecurrenceDaily
	case "weekly":
		return models.EventRecurrenceWeekly
	case "biweekly":
		return models.EventRecurrenceBiweekly
	case "monthly":
		return models.EventRecurrenceMonthly
	case "yearly":
		return models.EventRecurrenceYearly
	default:
		return models.EventRecurrenceNone
	}
}
