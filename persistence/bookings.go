package persistence

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// BookingRepository defines database operations for bookings.
type BookingRepository interface {
	// CreateBooking inserts a new booking record and returns the persisted model.
	CreateBooking(booking models.Booking) (*models.Booking, error)

	// GetBookingByID fetches a single booking by its UUID primary key.
	GetBookingByID(id string) (*models.Booking, error)

	// GetBookingsByTripID returns all bookings associated with the given trip.
	GetBookingsByTripID(tripID string) ([]models.Booking, error)
}

// CreateBooking inserts a new booking into the database.
func (r *RepositoryPg) CreateBooking(booking models.Booking) (*models.Booking, error) {
	if err := r.gormDB.Create(&booking).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to create booking: %w", err)
	}

	return &booking, nil
}

// GetBookingByID retrieves a booking by its primary key.
func (r *RepositoryPg) GetBookingByID(id string) (*models.Booking, error) {
	var booking models.Booking

	if err := r.gormDB.First(&booking, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to get booking with id %q: %w", id, err)
	}

	return &booking, nil
}

// GetBookingsByTripID returns all bookings for a trip ordered by creation time descending.
// Returns an empty slice (not an error) when the trip has no bookings.
func (r *RepositoryPg) GetBookingsByTripID(tripID string) ([]models.Booking, error) {
	var bookings []models.Booking

	if err := r.gormDB.
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to get bookings for trip %q: %w", tripID, err)
	}

	return bookings, nil
}