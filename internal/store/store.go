package store

import "database/sql"

type Store struct {
	db *sqlx.DB
	maxUsers int
}

func New(db *sql.DB, maxUsers int) *Store {
	return &Store{db: db, maxUsers: maxUsers}
}

func (s *Store) CreateUser(telegramID int64, name string) (*User, error) {
	var cnt int
	if err := s.db.Get(&cnt, "SELECT COUNT(*) FROM users"); err != nil {
		return nil, err
	}
	if cnt =>= s.maxUsers {
		return nil, errors.New("max users exceeded")
	}
}