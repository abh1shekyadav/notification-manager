package user

import "database/sql"

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) IsEmailRegistered(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresRepo) RegisterUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Email, user.Password, user.CreatedAt)
	return err
}

func (r *PostgresRepo) FindUserByEmail(email string) (*User, error) {
	row := r.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE email=$1", email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
