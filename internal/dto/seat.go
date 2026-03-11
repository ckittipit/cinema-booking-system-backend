package dto

type SeatResponse struct {
	SeatID    string  `json:"seat_id"`
	Status    string  `json:"status"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type SeatMapResponse struct {
	ShowtimeID string         `json:"showtime_id"`
	Seats      []SeatResponse `json:"seats"`
}
