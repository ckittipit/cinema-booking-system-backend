package service

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirebaseAuthService struct {
	app *firebase.App
}

func NewFirebaseAuthService(credentialsPath string) (*FirebaseAuthService, error) {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &FirebaseAuthService{
		app: app,
	}, nil
}

func (s *FirebaseAuthService) VerifyIDToken(ctx context.Context, idToken string) (map[string]any, error) {
	client, err := s.app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
