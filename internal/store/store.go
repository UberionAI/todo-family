package store

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db       *sqlx.DB
	maxUsers int
}

func New(db *sqlx.DB, maxUsers int) *Store {
	return &Store{db: db, maxUsers: maxUsers}
}

// USERS
func (s *Store) CreateUser(telegramID int64, name string) (*User, error) {
	var cnt int
	if err := s.db.Get(&cnt, "SELECT COUNT(*) FROM users"); err != nil {
		return nil, err
	}
	if cnt >= s.maxUsers {
		return nil, errors.New("max users exceeded")
	}

	id := uuid.New()
	u := &User{}
	err := s.db.Get(u,
		`INSERT INTO users (id, telegram_id, name)
		 VALUES ($1, $2, $3)
		 RETURNING id, telegram_id, name, created_at`,
		id, telegramID, name)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) ListUsers() ([]User, error) {
	var users []User
	if err := s.db.Select(&users,
		"SELECT id, telegram_id, name, created_at FROM users ORDER BY created_at"); err != nil {
		return nil, err
	}
	return users, nil
}

// GROUPS
func (s *Store) CreateGroup(name string, createdBy *uuid.UUID) (*Group, error) {
	id := uuid.New()
	g := &Group{}
	if err := s.db.Get(g,
		`INSERT INTO groups (id, name, created_by)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, created_by, created_at`,
		id, name, createdBy); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Store) RenameGroup(id uuid.UUID, newName string) error {
	res, err := s.db.Exec("UPDATE groups SET name = $1 WHERE id = $2", newName, id)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("group not found")
	}
	return nil
}

func (s *Store) DeleteGroup(id uuid.UUID) error {
	res, err := s.db.Exec("DELETE FROM groups WHERE id = $1", id)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("group not found")
	}
	return nil
}

func (s *Store) ListGroups() ([]Group, error) {
	var gs []Group
	if err := s.db.Select(&gs,
		"SELECT id, name, created_by, created_at FROM groups ORDER BY created_at"); err != nil {
		return nil, err
	}
	return gs, nil
}

// ENTRIES
func (s *Store) CreateEntry(groupID uuid.UUID, title string, description *string, createdBy *uuid.UUID) (*Entry, error) {
	id := uuid.New()
	e := &Entry{}
	if err := s.db.Get(e,
		`INSERT INTO entries (id, group_id, title, description, created_by)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, group_id, title, description, created_by, created_at, updated_at`,
		id, groupID, title, description, createdBy); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Store) RenameEntry(id uuid.UUID, newTitle string) error {
	res, err := s.db.Exec("UPDATE entries SET title = $1 WHERE id = $2", newTitle, id)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("entry not found")
	}
	return nil
}

func (s *Store) DeleteEntry(id uuid.UUID) error {
	res, err := s.db.Exec("DELETE FROM entries WHERE id = $1", id)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("entry not found")
	}
	return nil
}

func (s *Store) ListEntries(groupID uuid.UUID) ([]Entry, error) {
	var es []Entry
	if err := s.db.Select(&es,
		"SELECT id, group_id, title, description, created_by, created_at, updated_at FROM entries WHERE group_id = $1 ORDER BY created_at",
		groupID); err != nil {
		return nil, err
	}
	return es, nil
}
