package dto

type MovieResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"duration_minutes"`
	PosterURL       string `json:"poster_url"`
}

type ShowtimeResponse struct {
	ID          string `json:"id"`
	MovieID     string `json:"movie_id"`
	TheaterName string `json:"theater_name"`
	StartTime   string `'json:"start_time"`
	SeatRows    int    `json:"seat_rows"`
	SeatCols    int    `json:"seat_cols"`
}
