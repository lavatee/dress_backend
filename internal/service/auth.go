package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
)

type AuthService struct {
	repo *repository.Repository
}

const (
	salt       = "jd83420s32vv"
	tokenKey   = "xks7jhc93n94hz"
	accessTTL  = 15 * time.Minute
	refreshTTL = 15 * 24 * time.Hour
)

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateAdmin(name string, email string, password string) error {
	passwordHash := s.hashPassword(password)
	return s.repo.Auth.CreateAdmin(name, email, passwordHash)
}

func (s *AuthService) hashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	return fmt.Sprintf("%x", sha.Sum([]byte(salt)))
}

func (s *AuthService) SignUp(user model.User) (int, error) {
	user.Password = s.hashPassword(user.Password)
	return s.repo.Auth.CreateUser(user)
}

func (s *AuthService) SignIn(email, password string) (string, string, error) {
	passwordHash := s.hashPassword(password)
	userId, err := s.repo.Auth.SignIn(email, passwordHash)
	if err != nil {
		return "", "", err
	}
	accessClaims := jwt.MapClaims{
		"exp": time.Now().Add(accessTTL).Unix(),
		"id":  userId,
	}
	refreshClaims := jwt.MapClaims{
		"exp": time.Now().Add(refreshTTL).Unix(),
		"id":  userId,
	}
	access, err := s.NewToken(accessClaims)
	if err != nil {
		return "", "", err
	}
	refresh, err := s.NewToken(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *AuthService) NewToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		accessClaims := jwt.MapClaims{
			"exp": time.Now().Add(accessTTL).Unix(),
			"id":  claims["id"],
		}
		refreshClaims := jwt.MapClaims{
			"exp": time.Now().Add(refreshTTL).Unix(),
			"id":  claims["id"],
		}
		access, err := s.NewToken(accessClaims)
		if err != nil {
			return "", "", err
		}
		refresh, err := s.NewToken(refreshClaims)
		if err != nil {
			return "", "", err
		}
		return access, refresh, nil
	}
	return "", "", errors.New("invalid token")
}

func (s *AuthService) ParseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *AuthService) NewAdmin(thisAdminId, newAdminId int) error {
	if !s.repo.Auth.IsAdmin(thisAdminId) {
		return errors.New("user is not admin")
	}
	return s.repo.Auth.NewAdmin(thisAdminId, newAdminId)
}

func (s *AuthService) NewBuyer(thisAdminId, newBuyerId int) error {
	if !s.repo.Auth.IsAdmin(thisAdminId) {
		return errors.New("user is not admin")
	}
	return s.repo.Auth.NewBuyer(thisAdminId, newBuyerId)
}

func (s *AuthService) RemoveBuyer(thisAdminId, buyerId int) error {
	if !s.repo.Auth.IsAdmin(thisAdminId) {
		return errors.New("user is not admin")
	}
	return s.repo.Auth.RemoveBuyer(thisAdminId, buyerId)
}

func (s *AuthService) GetUserRole(userId int) (string, error) {
	return s.repo.Auth.GetUserRole(userId)
}

func (s *AuthService) GetUser(userId int) (model.User, error) {
	return s.repo.Auth.GetUser(userId)
}
