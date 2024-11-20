package services

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
)

type TokenService struct {
	repo repositories.Token
}

func NewTokenService(repo repositories.Token) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) Create(token models.Token) (int, error) {
	return s.repo.Create(token)
}

func (s *TokenService) GetByToken(token string) (*models.Token, error) {
	return s.repo.GetByToken(token)
}

func (s *TokenService) Update(token string) error {
	return s.repo.Update(token)
}

func (s *TokenService) UpdateYouTubeAccessToken(token string, userID int) error {
	return s.repo.UpdateYouTubeAccessToken(token, userID)
}

func (s *TokenService) RevokeAllUserTokens(userID int) error {
	return s.repo.RevokeAllUserTokens(userID)
}
