package store

import (
	"database/sql"
	"go_beginner/internals/tokens"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error 
	CrateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens ( hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (s *PostgresTokenStore) CrateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	newToken, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.Insert(newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func (s *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := s.db.Exec(query, userID, scope)
	return err
}
