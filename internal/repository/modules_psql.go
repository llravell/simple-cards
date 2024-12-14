package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/llravell/simple-cards/internal/entity"
)

type ModulesRepository struct {
	conn *sql.DB
}

func NewModulesRepository(conn *sql.DB) *ModulesRepository {
	return &ModulesRepository{conn: conn}
}

func (repo *ModulesRepository) GetAllModules(
	ctx context.Context,
	userUUID string,
) ([]*entity.Module, error) {
	modules := make([]*entity.Module, 0)

	rows, err := repo.conn.QueryContext(
		ctx,
		"SELECT uuid, name, user_uuid FROM modules WHERE user_uuid=$1",
		userUUID,
	)
	if err != nil {
		return modules, err
	}

	defer rows.Close()

	for rows.Next() {
		var module entity.Module

		err = rows.Scan(&module.UUID, &module.Name, &module.UserUUID)
		if err != nil {
			return modules, err
		}

		modules = append(modules, &module)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return modules, nil
		}

		return modules, err
	}

	return modules, nil
}

func (repo *ModulesRepository) CreateNewModule(
	ctx context.Context,
	userUUID string,
	moduleName string,
) (*entity.Module, error) {
	var module entity.Module

	row := repo.conn.QueryRowContext(ctx, `
		INSERT INTO modules (name, user_uuid)
		VALUES
			($1, $2)
		RETURNING uuid, name, user_uuid;
	`, moduleName, userUUID)

	err := row.Scan(&module.UUID, &module.Name, &module.UserUUID)
	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (repo *ModulesRepository) UpdateModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
	moduleName string,
) (*entity.Module, error) {
	var module entity.Module

	row := repo.conn.QueryRowContext(ctx, `
		UPDATE modules
		SET name=$1
		WHERE uuid=$2 AND user_uuid=$3
		RETURNING uuid, name, user_uuid;
	`, moduleName, moduleUUID, userUUID)

	err := row.Scan(&module.UUID, &module.Name, &module.UserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &entity.ModuleNotFoundError{UUID: moduleUUID}
		}

		return nil, err
	}

	return &module, nil
}

func (repo *ModulesRepository) DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error {
	_, err := repo.conn.ExecContext(ctx, `
		DELETE FROM modules
		WHERE uuid=$1 AND user_uuid=$2;
	`, moduleUUID, userUUID)

	return err
}
