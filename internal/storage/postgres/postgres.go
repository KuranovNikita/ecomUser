package postgres

import (
	"context"
	"database/sql"
	"ecomUser/internal/domain/models"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

const (
	driverName = "postgres"
)

type Storage struct {
	db *sql.DB
}

var ErrUserExists = errors.New("user with this email already exists")
var ErrUserNotFound = errors.New("user not found")

func New(storagePath string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open(driverName, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	storage := &Storage{
		db: db,
	}

	return storage, nil
}

func (s *Storage) Close() error {
	const op = "storage.Close"
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, login string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt := `INSERT INTO users(email, login, pass_hash) VALUES($1, $2, $3) RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, stmt, email, login, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w (email: %s)", op, ErrUserExists, email)
			}
		}
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.GetUser"

	stmt := `SELECT id, email, login, pass_hash FROM users WHERE id = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(&user.ID, &user.Email, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w (login: %s)", op, ErrUserNotFound, id)
		}
		return models.User{}, fmt.Errorf("%s: execute statement or scan: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserLogin(ctx context.Context, login string) (models.User, error) {
	const op = "storage.GetUser"

	stmt := `SELECT id, email, login, pass_hash FROM users WHERE login = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, stmt, login).Scan(&user.ID, &user.Email, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w (login: %s)", op, ErrUserNotFound, login)
		}
		return models.User{}, fmt.Errorf("%s: execute statement or scan: %w", op, err)
	}

	return user, nil
}

func SplitStoragePath(login, password, host, port, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", login, password, host, port, dbName)
}
