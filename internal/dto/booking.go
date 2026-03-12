package dto

type ConfirmBookingRequest struct {
	ShowtimeID string  `json:"showtime_id"`
	SeatID     string  `json:"seat_id"`
	Price      float64 `json:"price"`
}

type BookingResponse struct {
	ID         string  `json"id"`
	UserID     string  `json:"user_id"`
	ShowtimeID string  `json:"showtime_id"`
	SeatID     string  `json:"seat_id"`
	Status     string  `json:"status`
	Price      float64 `json:"price"`
	CreatedAt  string  `json:"created_at"`
}
