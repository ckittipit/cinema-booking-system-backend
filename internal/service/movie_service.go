package service

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/repository"
	"context"
)

type MovieService struct {
	movieRepository *repository.MovieRepository
}

func NewMovieService(movieRepository *repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepository: movieRepository,
	}
}

//	func (s *MovieService) GetMovies(ctx context.Context) (any, error) {
//		// ใข้ any ไปก่อน พอเริ่มคล่องจะกลับมาเปลี่ยนtypeให้สวยขึ้น
//		return s.movieRepository.FindAll(ctx)
//	}
func (s *MovieService) GetMovies(ctx context.Context) ([]dto.MovieResponse, error) {
	movies, err := s.movieRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MovieResponse, 0, len(movies))
	for _, movie := range movies {
		result = append(result, dto.MovieResponse{
			ID:              movie.ID.Hex(),
			Title:           movie.Title,
			Description:     movie.Description,
			DurationMinutes: movie.DurationMinutes,
			PosterURL:       movie.PosterURL,
		})
	}

	return result, nil
}
