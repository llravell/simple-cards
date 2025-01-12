package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var module entity.Module

		err = rows.Scan(&module.UUID, &module.Name, &module.UserUUID)
		if err != nil {
			return nil, err
		}

		modules = append(modules, &module)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return modules, nil
		}

		return nil, err
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

func (repo *ModulesRepository) CreateNewModuleWithCards(
	ctx context.Context,
	moduleWithCards *entity.ModuleWithCards,
) error {
	tx, err := repo.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO modules (name, user_uuid)
		VALUES ($1, $2)
		RETURNING uuid;
	`, moduleWithCards.Name, moduleWithCards.UserUUID)

	err = row.Scan(&moduleWithCards.UUID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	insertColsAmount := 3
	insertParts := make([]string, 0, len(moduleWithCards.Cards))
	args := make([]any, 0, len(moduleWithCards.Cards)*insertColsAmount)

	for i, card := range moduleWithCards.Cards {
		base := i * insertColsAmount
		part := fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3)

		insertParts = append(insertParts, part)
		args = append(args, moduleWithCards.UUID, card.Term, card.Meaning)
	}

	//nolint:gosec
	query := fmt.Sprintf(`
		INSERT INTO cards (module_uuid, term, meaning)
		VALUES %s;
	`, strings.Join(insertParts, ","))

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
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

func (repo *ModulesRepository) DeleteModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) error {
	_, err := repo.conn.ExecContext(ctx, `
		DELETE FROM modules
		WHERE uuid=$1 AND user_uuid=$2;
	`, moduleUUID, userUUID)

	return err
}

func (repo *ModulesRepository) GetModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) (*entity.Module, error) {
	var module entity.Module

	row := repo.conn.QueryRowContext(ctx, `
		SELECT uuid, name, user_uuid
		FROM modules
		WHERE uuid=$1 AND user_uuid=$2;
	`, moduleUUID, userUUID)

	err := row.Scan(&module.UUID, &module.Name, &module.UserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &entity.ModuleNotFoundError{UUID: moduleUUID}
		}

		return nil, err
	}

	return &module, nil
}

func (repo *ModulesRepository) ModuleExists(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) (bool, error) {
	row := repo.conn.QueryRowContext(ctx, `
		SELECT uuid
		FROM modules
		WHERE uuid=$1 AND user_uuid=$2;
	`, moduleUUID, userUUID)

	err := row.Scan(&moduleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
