package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/llravell/simple-cards/internal/entity"
)

type UsersRepository struct {
	conn *sql.DB
}

func NewUsersRepository(conn *sql.DB) *UsersRepository {
	return &UsersRepository{conn: conn}
}

func (r *UsersRepository) StoreUser(
	ctx context.Context,
	login string,
	passwordHash string,
) (*entity.User, error) {
	var user entity.User

	row := r.conn.QueryRowContext(ctx, `
		INSERT INTO users (login, password)
		VALUES
			($1, $2)
		RETURNING uuid, login, password;
	`, login, passwordHash)

	err := row.Scan(&user.UUID, &user.Login, &user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = entity.ErrUserConflict
		}
	}

	return &user, err
}

func (r *UsersRepository) FindUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User

	row := r.conn.QueryRowContext(ctx, `
		SELECT uuid, login, password
		FROM users
		WHERE
			login=$1;
	`, login)

	err := row.Scan(&user.UUID, &user.Login, &user.Password)

	return &user, err
}
