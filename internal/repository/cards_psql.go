package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (repo *CardsRepository) CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	var storedCard entity.Card

	row := repo.conn.QueryRowContext(ctx, `
		INSERT INTO cards (module_uuid, term, meaning)
		VALUES
			($1, $2, $3)
		RETURNING uuid, module_uuid, term, meaning;
	`, card.ModuleUUID, card.Term, card.Meaning)

	err := row.Scan(&storedCard.UUID, &storedCard.ModuleUUID, &storedCard.Term, &storedCard.Meaning)

	return &storedCard, err
}

func (repo *CardsRepository) SaveCard(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	var storedCard entity.Card

	updatedFields := make([]string, 0)
	setParts := make([]string, 0)
	args := make([]any, 0)

	if card.Term != "" {
		updatedFields = append(updatedFields, "term")
		args = append(args, card.Term)
	}

	if card.Meaning != "" {
		updatedFields = append(updatedFields, "meaning")
		args = append(args, card.Meaning)
	}

	for i, filed := range updatedFields {
		part := fmt.Sprintf("%s=$%d", filed, i+1)
		setParts = append(setParts, part)
	}

	//nolint:gosec
	query := fmt.Sprintf(`
		UPDATE cards
		SET %s
		WHERE uuid=$%d AND module_uuid=$%d
		RETURNING uuid, module_uuid, term, meaning;
	`, strings.Join(setParts, ","), len(setParts)+1, len(setParts)+2)

	args = append(args, card.UUID, card.ModuleUUID)
	row := repo.conn.QueryRowContext(ctx, query, args...)
	err := row.Scan(&storedCard.UUID, &storedCard.ModuleUUID, &storedCard.Term, &storedCard.Meaning)

	return &storedCard, err
}

func (repo *CardsRepository) DeleteCard(ctx context.Context, moduleUUID string, cardUUID string) error {
	_, err := repo.conn.ExecContext(ctx, `
		DELETE FROM cards
		WHERE uuid=$1 AND module_uuid=$2;
	`, cardUUID, moduleUUID)

	return err
}
