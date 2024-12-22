package usecase

import (
	"context"

	"github.com/llravell/simple-cards/internal/entity"
)

type ModulesUseCase struct {
	modulesRepo         ModulesRepository
	cardsRepo           CardsRepository
	quizletModuleParser QuizletModuleParser
}

func NewModulesUseCase(
	modulesRepo ModulesRepository,
	cardsRepo CardsRepository,
	quizletModuleParser QuizletModuleParser,
) *ModulesUseCase {
	return &ModulesUseCase{
		modulesRepo:         modulesRepo,
		cardsRepo:           cardsRepo,
		quizletModuleParser: quizletModuleParser,
	}
}

func (uc *ModulesUseCase) GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error) {
	return uc.modulesRepo.GetAllModules(ctx, userUUID)
}

func (uc *ModulesUseCase) ModuleExists(ctx context.Context, userUUID string, moduleUUID string) (bool, error) {
	return uc.modulesRepo.ModuleExists(ctx, userUUID, moduleUUID)
}

func (uc *ModulesUseCase) CreateNewModule(
	ctx context.Context,
	userUUID string,
	moduleName string,
) (*entity.Module, error) {
	return uc.modulesRepo.CreateNewModule(ctx, userUUID, moduleName)
}

func (uc *ModulesUseCase) UpdateModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
	moduleName string,
) (*entity.Module, error) {
	return uc.modulesRepo.UpdateModule(ctx, userUUID, moduleUUID, moduleName)
}

func (uc *ModulesUseCase) DeleteModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) error {
	return uc.modulesRepo.DeleteModule(ctx, userUUID, moduleUUID)
}

func (uc *ModulesUseCase) GetModuleWithCards(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) (*entity.ModuleWithCards, error) {
	module, err := uc.modulesRepo.GetModule(ctx, userUUID, moduleUUID)
	if err != nil {
		return nil, err
	}

	cards, err := uc.cardsRepo.GetModuleCards(ctx, moduleUUID)
	if err != nil {
		return nil, err
	}

	return &entity.ModuleWithCards{
		Module: *module,
		Cards:  cards,
	}, nil
}

func (uc *ModulesUseCase) ImportModuleFromQuizlet(
	ctx context.Context,
	userUUID string,
	moduleName string,
	quizletModuleID string,
) {
	quizletCards, err := uc.quizletModuleParser.Parse(quizletModuleID)
	if err != nil {
		return
	}

	if len(quizletCards) == 0 {
		return
	}

	module := entity.Module{
		Name:     moduleName,
		UserUUID: userUUID,
	}

	moduleCards := make([]*entity.Card, 0, len(quizletCards))

	for _, quizletCard := range quizletCards {
		card := &entity.Card{
			Term:    quizletCard.Front,
			Meaning: quizletCard.Back,
		}

		moduleCards = append(moduleCards, card)
	}

	uc.modulesRepo.CreateNewModuleWithCards(
		ctx,
		&entity.ModuleWithCards{
			Module: module,
			Cards:  moduleCards,
		},
	)
}
