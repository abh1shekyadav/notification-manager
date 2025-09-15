package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput           = errors.New("invalid input")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, ErrInvalidInput
	}
	registered, err := s.repo.IsEmailRegistered(email)
	if err != nil {
		return nil, err
	}
	if registered {
		return nil, ErrEmailAlreadyRegistered
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        uuid.NewString,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	if err := s.repo.RegisterUser(user); err != nil {
		return nil, err
	}
	return user, nil

}
