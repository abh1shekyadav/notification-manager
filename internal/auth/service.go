package auth

import (
	"errors"

	"github.com/abh1shekyadav/notification-manager/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  user.UserRepository
	validator AuthValidator
	secret    string
}

func NewAuthService(userRepo user.UserRepository, secret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		secret:   secret,
	}
}

func (s *AuthService) Login(email, password, secret string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}
	token, err := GenerateJWT(user.ID, user.Email, secret)
	if err != nil {
		return "", err
	}
	return token, nil
}
