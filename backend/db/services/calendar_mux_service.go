package services

import (
	"errors"
	"family-calendar-backend/db"
	"family-calendar-backend/db/models"
)

// CreateCalendarMux creates a new calendar mux for a user
func CreateCalendarMux(userID uint, name, description string) (*models.CalendarMux, error) {
	calendarMux := &models.CalendarMux{
		CreatedByID: userID,
		Name:        name,
		Description: description,
	}

	result := db.DB.Create(calendarMux)
	if result.Error != nil {
		return nil, result.Error
	}

	return calendarMux, nil
}

// GetCalendarMuxesByUser returns all calendar muxes created by a specific user
func GetCalendarMuxesByUser(userID uint) ([]models.CalendarMux, error) {
	var calendarMuxes []models.CalendarMux
	result := db.DB.Where("created_by_id = ?", userID).Find(&calendarMuxes)
	if result.Error != nil {
		return nil, result.Error
	}

	return calendarMuxes, nil
}

// DeleteCalendarMux deletes a calendar mux if it belongs to the specified user
func DeleteCalendarMux(id, userID uint) error {
	result := db.DB.Where("id = ? AND created_by_id = ?", id, userID).Delete(&models.CalendarMux{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("calendar mux not found or access denied")
	}

	return nil
}
