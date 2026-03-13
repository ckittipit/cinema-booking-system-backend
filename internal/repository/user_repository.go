package repository

import (
	"cinema-booking/backend/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	colloction *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		colloction: db.Collection("users"),
	}
}

func (r *UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user model.User
	if err := r.colloction.FindOne(ctx, bson.M{"firebase_uid": firebaseUID}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
