package auth

import (
	"errors"

	"github.com/abh1shekyadav/notification-manager/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  user.UserRepository
	validator AuthValidator
}

func NewAuthService(userRepo user.UserRepository, validator AuthValidator) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}
	token, err := GenerateJWT(user.ID, user.Email, "supersecret")
	if err != nil {
		return "", err
	}
	return token, nil
}
