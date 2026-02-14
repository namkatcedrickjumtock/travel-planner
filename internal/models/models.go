package models

import "time"

// TripStatus represents the lifecycle state of a trip.
type TripStatus string

const (
	TripStatusPlanning  TripStatus = "planning"
	TripStatusConfirmed TripStatus = "confirmed"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

// BookingType identifies what kind of item is being booked.
type BookingType string

const (
	BookingTypeHotel    BookingType = "hotel"
	BookingTypeFlight   BookingType = "flight"
	BookingTypeActivity BookingType = "activity"
)

// BookingStatus represents the state of a booking.
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Trip is a top-level travel itinerary owned by a user.
type Trip struct {
	ID          string     `json:"id"           gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string     `json:"user_id"      gorm:"type:uuid;not null;index"`
	Title       string     `json:"title"        gorm:"type:varchar;not null"`
	Destination string     `json:"destination"  gorm:"type:varchar;not null"`
	StartDate   time.Time  `json:"start_date"   gorm:"not null"`
	EndDate     time.Time  `json:"end_date"     gorm:"not null"`
	Status      TripStatus `json:"status"       gorm:"type:varchar;not null;default:'planning'"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Hotel represents an accommodation option available for booking.
type Hotel struct {
	ID            string    `json:"id"             gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string    `json:"name"           gorm:"type:varchar;not null"`
	Location      string    `json:"location"       gorm:"type:varchar;not null"`
	PricePerNight float64   `json:"price_per_night" gorm:"type:numeric(10,2);not null"`
	Rating        float64   `json:"rating"         gorm:"type:numeric(3,2)"`
	AvailableFrom time.Time `json:"available_from"`
	AvailableTo   time.Time `json:"available_to"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Flight represents a flight available for booking.
type Flight struct {
	ID             string    `json:"id"               gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Airline        string    `json:"airline"          gorm:"type:varchar;not null"`
	Origin         string    `json:"origin"           gorm:"type:varchar(100);not null"`
	Destination    string    `json:"destination"      gorm:"type:varchar(100);not null"`
	DepartureTime  time.Time `json:"departure_time"   gorm:"not null"`
	ArrivalTime    time.Time `json:"arrival_time"     gorm:"not null"`
	Price          float64   `json:"price"            gorm:"type:numeric(10,2);not null"`
	SeatsAvailable int       `json:"seats_available"  gorm:"not null;default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Activity represents a bookable attraction or experience at a destination.
type Activity struct {
	ID            string    `json:"id"             gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string    `json:"name"           gorm:"type:varchar;not null"`
	Location      string    `json:"location"       gorm:"type:varchar;not null"`
	Description   string    `json:"description"    gorm:"type:text"`
	Price         float64   `json:"price"          gorm:"type:numeric(10,2);not null"`
	DurationHours float64   `json:"duration_hours" gorm:"type:numeric(5,2);not null"`
	AvailableDate time.Time `json:"available_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Booking links a Trip to a Hotel, Flight, or Activity.
// ReferenceID points to the ID of the booked item (hotel/flight/activity).
type Booking struct {
	ID          string        `json:"id"           gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TripID      string        `json:"trip_id"      gorm:"type:uuid;not null;index"`
	Type        BookingType   `json:"type"         gorm:"type:varchar;not null"`
	ReferenceID string        `json:"reference_id" gorm:"type:uuid;not null"`
	Status      BookingStatus `json:"status"       gorm:"type:varchar;not null;default:'pending'"`
	TotalPrice  float64       `json:"total_price"  gorm:"type:numeric(10,2);not null"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ─────────────────────────────────────────────
// Request / Response DTOs
// ─────────────────────────────────────────────

// CreateTripRequest is the payload for creating a new trip.
type CreateTripRequest struct {
	UserID      string    `json:"user_id"     binding:"required"`
	Title       string    `json:"title"       binding:"required"`
	Destination string    `json:"destination" binding:"required"`
	StartDate   time.Time `json:"start_date"  binding:"required"`
	EndDate     time.Time `json:"end_date"    binding:"required"`
}

// UpdateTripRequest is the payload for updating an existing trip.
// All fields are optional — only provided fields will be updated.
type UpdateTripRequest struct {
	Title       *string     `json:"title"`
	Destination *string     `json:"destination"`
	StartDate   *time.Time  `json:"start_date"`
	EndDate     *time.Time  `json:"end_date"`
	Status      *TripStatus `json:"status"`
}

// TripSearchParams carries filter criteria for searching trips.
type TripSearchParams struct {
	Destination string    `form:"destination"`
	StartDate   time.Time `form:"start_date"  time_format:"2006-01-02"`
	EndDate     time.Time `form:"end_date"    time_format:"2006-01-02"`
}

// CreateBookingRequest is the payload for booking an item within a trip.
type CreateBookingRequest struct {
	Type        BookingType `json:"type"         binding:"required"`
	ReferenceID string      `json:"reference_id" binding:"required"`
	TotalPrice  float64     `json:"total_price"  binding:"required,gt=0"`
}

// ErrorResponse is a uniform error envelope returned by all endpoints.
type ErrorResponse struct {
	Error string `json:"error"`
}
