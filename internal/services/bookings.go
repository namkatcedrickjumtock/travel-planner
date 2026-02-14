package services

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// BookingService defines business operations for bookings.
type BookingService interface {
	// BookItem creates a booking linking a trip to a hotel, flight, or activity.
	BookItem(tripID string, req models.CreateBookingRequest) (*models.Booking, error)

	// GetBooking retrieves a single booking by ID.
	GetBooking(id string) (*models.Booking, error)

	// GetTripBookings returns all bookings associated with the given trip.
	GetTripBookings(tripID string) ([]models.Booking, error)
}

// BookItem validates the request, verifies both the trip and the referenced
// item exist, then persists the booking.
func (s *TravelPlannerServiceImpl) BookItem(tripID string, req models.CreateBookingRequest) (*models.Booking, error) {
	if tripID == "" {
		return nil, fmt.Errorf("services: trip_id must not be empty")
	}

	// Verify the parent trip exists before creating a booking against it.
	if _, err := s.repo.GetTripByID(tripID); err != nil {
		return nil, fmt.Errorf("services: trip not found for booking: %w", err)
	}

	// Verify the referenced item exists to prevent orphaned bookings.
	if err := s.verifyReferenceExists(req.Type, req.ReferenceID); err != nil {
		return nil, err
	}

	booking := models.Booking{
		TripID:      tripID,
		Type:        req.Type,
		ReferenceID: req.ReferenceID,
		Status:      models.BookingStatusPending,
		TotalPrice:  req.TotalPrice,
	}

	created, err := s.repo.CreateBooking(booking)
	if err != nil {
		return nil, fmt.Errorf("services: create booking failed: %w", err)
	}

	return created, nil
}

// GetBooking retrieves a booking by its UUID.
func (s *TravelPlannerServiceImpl) GetBooking(id string) (*models.Booking, error) {
	if id == "" {
		return nil, fmt.Errorf("services: booking id must not be empty")
	}

	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return nil, fmt.Errorf("services: get booking failed: %w", err)
	}

	return booking, nil
}

// GetTripBookings returns all bookings for a trip.
func (s *TravelPlannerServiceImpl) GetTripBookings(tripID string) ([]models.Booking, error) {
	if tripID == "" {
		return nil, fmt.Errorf("services: trip_id must not be empty")
	}

	// Verify the trip exists so we return a 404 rather than an empty list
	// when the caller provides an unknown trip ID.
	if _, err := s.repo.GetTripByID(tripID); err != nil {
		return nil, fmt.Errorf("services: trip not found: %w", err)
	}

	bookings, err := s.repo.GetBookingsByTripID(tripID)
	if err != nil {
		return nil, fmt.Errorf("services: get trip bookings failed: %w", err)
	}

	return bookings, nil
}

// verifyReferenceExists checks that the item being booked actually exists in
// the database, routing to the correct repository method by booking type.
func (s *TravelPlannerServiceImpl) verifyReferenceExists(bookingType models.BookingType, referenceID string) error {
	var err error

	switch bookingType {
	case models.BookingTypeHotel:
		_, err = s.repo.GetHotelByID(referenceID)
	case models.BookingTypeFlight:
		_, err = s.repo.GetFlightByID(referenceID)
	default:
		return fmt.Errorf("services: unsupported booking type %q", bookingType)
	}

	if err != nil {
		return fmt.Errorf("services: referenced %s with id %q not found: %w", bookingType, referenceID, err)
	}

	return nil
}