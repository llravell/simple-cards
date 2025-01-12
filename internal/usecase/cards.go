package usecase

import (
	"context"

	"github.com/llravell/simple-cards/internal/entity"
)

type CardsUseCase struct {
	repo CardsRepository
}

func NewCardsUseCase(cardsRepo CardsRepository) *CardsUseCase {
	return &CardsUseCase{
		repo: cardsRepo,
	}
}

func (uc *CardsUseCase) GetModuleCards(ctx context.Context, moduleUUID string) ([]*entity.Card, error) {
	return uc.repo.GetModuleCards(ctx, moduleUUID)
}

func (uc *CardsUseCase) CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	return uc.repo.CreateCard(ctx, card)
}

func (uc *CardsUseCase) SaveCard(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	return uc.repo.SaveCard(ctx, card)
}

func (uc *CardsUseCase) DeleteCard(ctx context.Context, moduleUUID string, cardUUID string) error {
	return uc.repo.DeleteCard(ctx, moduleUUID, cardUUID)
}
