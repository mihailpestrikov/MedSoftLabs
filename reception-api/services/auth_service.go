package services

import (
	"errors"
	"reception-api/database"
	"reception-api/middleware"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *database.Repository
	jwtService *middleware.JWTService
}

func NewAuthService(repo *database.Repository, jwtService *middleware.JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(username, password string) (accessToken, refreshToken string, err error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err = s.jwtService.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(username, string(hash))
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
