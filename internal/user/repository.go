package user

type UserRepository interface {
	IsEmailRegistered(email string) (bool, error)
	RegisterUser(user *User) error
	// FindUserByEmail(email string) (*User, error)
}

type InMemoryRepo struct {
	users map[string]User
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{users: make(map[string]User)}
}

func (r *InMemoryRepo) IsEmailRegistered(email string) (bool, error) {
	_, exists := r.users[email]
	return exists, nil
}

func (r *InMemoryRepo) RegisterUser(u *User) error {
	r.users[u.Email] = *u
	return nil
}
