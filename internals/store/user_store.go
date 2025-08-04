package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainText
	p.hash = hash
	return nil
}

func (p *password) Check(plainText string) (bool, error) {
	if plainText == "" {
		return false , errors.New("password is not set")
	}
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText)) 
	if err != nil {
		return false, err
	}
	return true, nil
}

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(id int64, user *User) error
	DeleteUser(id int64) error
	GetUsers() ([]User, error)
	GetUserByEmail(email string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `INSERT INTO users (name, email, password, bio) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(query, user.Name, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	query := `SELECT id, name, email, password, bio, created_at, updated_at FROM users WHERE id = $1`
	var user User
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) UpdateUser(id int64, user *User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3, bio = $4, updated_at = NOW() WHERE id = $5`
	_, err := s.db.Exec(query, user.Name, user.Email, user.PasswordHash.hash, user.Bio, id)
	return err
}
func (s *PostgresUserStore) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
func (s *PostgresUserStore) GetUsers() ([]User, error) {
	query := `SELECT id, name, email, password, bio, created_at, updated_at FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password, bio, created_at, updated_at FROM users WHERE email = $1`
	var user User
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}