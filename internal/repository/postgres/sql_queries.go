package postgres

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
)

func (p *PostgresStorage) CreateUser(ctx context.Context, login, password string) error {
	query := `INSERT INTO users (login, password) VALUES ($1, $2)`
	_, err := p.DB.ExecContext(ctx, query, login, password)
	return err
}

func (p *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (*auth.User, error) {
	query := `SELECT id, login, password FROM users WHERE login = $1`
	row := p.DB.QueryRowContext(ctx, query, login)
	user := &auth.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	return user, err
}
