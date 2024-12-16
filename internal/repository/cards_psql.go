package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/llravell/simple-cards/internal/entity"
)

type CardsRepository struct {
	conn *sql.DB
}

func NewCardsRepository(conn *sql.DB) *CardsRepository {
	return &CardsRepository{conn: conn}
}

func (repo *CardsRepository) GetModuleCards(
	ctx context.Context,
	moduleUUID string,
) ([]*entity.Card, error) {
	cards := make([]*entity.Card, 0)

	rows, err := repo.conn.QueryContext(ctx, `
		SELECT uuid, term, meaning, module_uuid
		FROM cards
		WHERE module_uuid=$1;
	`, moduleUUID)
	if err != nil {
		return cards, err
	}

	defer rows.Close()

	for rows.Next() {
		var card entity.Card

		err = rows.Scan(&card.UUID, &card.Term, &card.Meaning, &card.ModuleUUID)
		if err != nil {
			return cards, err
		}

		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cards, nil
		}

		return cards, err
	}

	return cards, nil
}
