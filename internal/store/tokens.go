package store

import (
	"database/sql"
	"time"

	"github.com/DavidGudovic/api_exercise/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, tokenType string, ttlMinutes time.Duration) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, tokenType string) error
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
			  VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (s *PostgresTokenStore) CreateNewToken(userID int, tokenType string, ttlMinutes time.Duration) (*tokens.Token, error) {
	ttl := ttlMinutes * time.Minute
	token, err := tokens.GenerateToken(userID, ttl, tokenType)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *PostgresTokenStore) DeleteAllTokensForUser(userID int, tokenType string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	_, err := s.db.Exec(query, userID, tokenType)
	return err
}
