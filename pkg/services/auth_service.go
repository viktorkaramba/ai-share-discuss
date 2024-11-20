package services

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const (
	tokenTTL   = 12 * time.Hour
	signingKey = "fewgf233io4y9238y0h239g23"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId       int    `json: "user_id"`
	UserPlatform string `json: "userPlatform"`
}

type AuthService struct {
	repoAuth  repositories.Authorization
	repoToken repositories.Token
}

func NewAuthService(repoAuth repositories.Authorization, repoToken repositories.Token) *AuthService {
	return &AuthService{repoAuth: repoAuth, repoToken: repoToken}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	return s.repoAuth.CreateUser(user)
}

func (s *AuthService) CreateUserSpotify(spotifyUri string, token models.ApiToken, user models.User) (int, error) {
	return s.repoAuth.CreateUserSpotify(spotifyUri, token, user)
}

func (s *AuthService) CreateUserYouTubeMusic(token string, user models.User) (int, error) {
	return s.repoAuth.CreateUserYouTubeMusic(token, user)
}

func (s *AuthService) GetUser(email string) (*models.User, error) {
	return s.repoAuth.GetUser(email)
}

func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.repoAuth.GetUserByID(id)
}

func (s *AuthService) GenerateToken(email string) (string, error) {
	user, err := s.repoAuth.GetUser(email)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:       user.ID,
		UserPlatform: user.Platform,
	})
	return token.SignedString([]byte(os.Getenv("SIGNINKEY")))
}

func (s *AuthService) RefreshToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
	})
	tokenSigned, err := token.SignedString([]byte(signingKey))
	err = s.repoToken.RevokeAllUserTokens(userId)
	if err != nil {
		return "", err
	}
	return tokenSigned, err
}

func (s *AuthService) ParseToken(accessToken string) (*models.UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SIGNINKEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	userClaims := models.UserClaims{UserID: claims.UserId, UserPlatform: claims.UserPlatform}
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}
	return &userClaims, nil
}
