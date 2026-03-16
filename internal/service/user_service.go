package service

import (
	"cinema-booking/backend/internal/model"
	"cinema-booking/backend/internal/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) FindOrCreateUser(
	ctx context.Context,
	firebaseUID string,
	email string,
	name string,
) (*model.User, error) {
	user, err := s.userRepository.FindByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		return user, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	now := time.Now()
	newUser := &model.User{
		FirebaseUID: firebaseUID,
		Email:       email,
		Name:        name,
		Role:        model.UserRoleUser,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.userRepository.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
