package user

import (
	"errors"
	"log"
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
		log.Println("RegisterUser error: empty email or password")
		return nil, ErrInvalidInput
	}

	registered, err := s.repo.IsEmailRegistered(email)
	if err != nil {
		log.Println("RegisterUser error checking if email is registered:", err)
		return nil, err
	}
	if registered {
		log.Println("RegisterUser error: email already registered:", email)
		return nil, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("RegisterUser error hashing password:", err)
		return nil, err
	}

	user := &User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.repo.RegisterUser(user); err != nil {
		log.Println("RegisterUser error inserting user into repository:", err)
		return nil, err
	}

	log.Println("RegisterUser success:", email)
	return user, nil
}

func (s *UserService) FindUserByEmail(email string) (*User, error) {
	if email == "" {
		log.Println("FindUserByEmail error: empty email")
		return nil, ErrInvalidInput
	}

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		log.Println("FindUserByEmail repository error:", err)
		return nil, err
	}

	log.Println("FindUserByEmail success:", email)
	return user, nil
}
