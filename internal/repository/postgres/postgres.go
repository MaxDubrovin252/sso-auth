package postgres

import (
	"context"
	"fmt"
	"my-grpc/internal/config"
	"my-grpc/internal/lib/models"
	"os"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg config.DBConfig) (*Storage, error) {
	const op = "repos.postgres.newDB"
	password := os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Port, cfg.Host, cfg.UserName, password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, password []byte) (int64, error) {
	const op = "repos.postgres.saveuser"

	stmt, err := s.db.Prepare("INSERT INTO users (email,pass_hash) VALUES($1,$2)")

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, password)

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "repos.postgres.user"

	stmt, err := s.db.Prepare("SELECT id,email,pass_hash FROM users WHERE email = $1")

	if err != nil {
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user models.User

	if err := row.Scan(&user.Id, &user.Email, &user.PassHash); err != nil {
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "repos.postgres.isadmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = $1")

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return isAdmin, nil
}
func (s *Storage) App(ctx context.Context, appId int) (models.App, error) {
	panic("asd")
}
