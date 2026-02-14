package persistence

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// TripRepository defines all database operations for trips.
type TripRepository interface {
	// CreateTrip inserts a new trip record and returns the persisted model.
	CreateTrip(trip models.Trip) (*models.Trip, error)

	// GetTripByID fetches a single trip by its UUID primary key.
	GetTripByID(id string) (*models.Trip, error)

	// GetAllTrips returns every trip in the database.
	GetAllTrips() ([]models.Trip, error)

	// UpdateTrip applies a partial update map to the trip with the given ID.
	// Only the keys present in `updates` are written to the database.
	UpdateTrip(id string, updates map[string]interface{}) (*models.Trip, error)

	// DeleteTrip hard-deletes the trip with the given ID.
	DeleteTrip(id string) error

	// SearchTrips filters trips by destination and/or date range.
	// Any zero-value filter field is ignored.
	SearchTrips(params models.TripSearchParams) ([]models.Trip, error)
}

// CreateTrip inserts a new trip into the database.
func (r *RepositoryPg) CreateTrip(trip models.Trip) (*models.Trip, error) {
	if err := r.gormDB.Create(&trip).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to create trip: %w", err)
	}

	return &trip, nil
}

// GetTripByID retrieves a trip by its primary key.
// Returns an error wrapping gorm.ErrRecordNotFound when no row matches.
func (r *RepositoryPg) GetTripByID(id string) (*models.Trip, error) {
	var trip models.Trip

	if err := r.gormDB.First(&trip, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to get trip with id %q: %w", id, err)
	}

	return &trip, nil
}

// GetAllTrips returns all trips ordered by creation time descending.
func (r *RepositoryPg) GetAllTrips() ([]models.Trip, error) {
	var trips []models.Trip

	if err := r.gormDB.Order("created_at DESC").Find(&trips).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to list trips: %w", err)
	}

	return trips, nil
}

// UpdateTrip applies the provided field map to the trip row and returns the
// updated record. An empty updates map is a no-op but not an error.
func (r *RepositoryPg) UpdateTrip(id string, updates map[string]interface{}) (*models.Trip, error) {
	// Confirm the trip exists before attempting to update.
	trip, err := r.GetTripByID(id)
	if err != nil {
		return nil, fmt.Errorf("persistence: update pre-check failed: %w", err)
	}

	if err := r.gormDB.Model(trip).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to update trip with id %q: %w", id, err)
	}

	return trip, nil
}

// DeleteTrip removes the trip row with the given ID.
// Returns an error if no row was deleted (i.e. ID not found).
func (r *RepositoryPg) DeleteTrip(id string) error {
	result := r.gormDB.Delete(&models.Trip{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("persistence: failed to delete trip with id %q: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("persistence: no trip found with id %q", id)
	}

	return nil
}

// SearchTrips applies the non-zero fields in params as WHERE filters.
// destination is a case-insensitive partial match; date fields are range filters.
func (r *RepositoryPg) SearchTrips(params models.TripSearchParams) ([]models.Trip, error) {
	query := r.gormDB.Model(&models.Trip{})

	if params.Destination != "" {
		// ILIKE enables case-insensitive partial matching on destination.
		query = query.Where("destination ILIKE ?", "%"+params.Destination+"%")
	}

	if !params.StartDate.IsZero() {
		query = query.Where("start_date >= ?", params.StartDate)
	}

	if !params.EndDate.IsZero() {
		query = query.Where("end_date <= ?", params.EndDate)
	}

	var trips []models.Trip
	if err := query.Order("start_date ASC").Find(&trips).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to search trips: %w", err)
	}

	return trips, nil
}