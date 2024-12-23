package usecase

import (
	"context"

	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

//go:generate ../../bin/mockgen -source=modules.go -destination=../mocks/mock_modules_usecase.go -package=mocks
type QuizletImportWorkerPool interface {
	QueueWork(w *QuizletImportWork) error
}

type QuizletImportWork struct {
	repo                ModulesRepository
	quizletModuleParser QuizletModuleParser
	log                 *zerolog.Logger
	module              *entity.Module
	quizletModuleID     string
}

func (w *QuizletImportWork) Do(ctx context.Context) {
	quizletCards, err := w.quizletModuleParser.Parse(ctx, w.quizletModuleID)
	if err != nil {
		w.log.Error().Err(err).Msg("quizlet module parsing failed")

		return
	}

	if len(quizletCards) == 0 {
		return
	}

	w.log.Info().Msgf("quizlet module \"%s\" parsed", w.quizletModuleID)

	moduleCards := make([]*entity.Card, 0, len(quizletCards))

	for _, quizletCard := range quizletCards {
		card := &entity.Card{
			Term:    quizletCard.Front,
			Meaning: quizletCard.Back,
		}

		moduleCards = append(moduleCards, card)
	}

	err = w.repo.CreateNewModuleWithCards(
		ctx,
		&entity.ModuleWithCards{
			Module: *w.module,
			Cards:  moduleCards,
		},
	)
	if err != nil {
		w.log.Error().Err(err).Msg("module from quizlet storing failed")
	} else {
		w.log.Info().Msgf("quizlet module \"%s\" imported", w.quizletModuleID)
	}
}

type ModulesUseCase struct {
	modulesRepo         ModulesRepository
	cardsRepo           CardsRepository
	quizletModuleParser QuizletModuleParser
	quizletImportWP     QuizletImportWorkerPool
	log                 *zerolog.Logger
}

func NewModulesUseCase(
	modulesRepo ModulesRepository,
	cardsRepo CardsRepository,
	quizletModuleParser QuizletModuleParser,
	quizletImportWP QuizletImportWorkerPool,
	log *zerolog.Logger,
) *ModulesUseCase {
	return &ModulesUseCase{
		modulesRepo:         modulesRepo,
		cardsRepo:           cardsRepo,
		quizletModuleParser: quizletModuleParser,
		quizletImportWP:     quizletImportWP,
		log:                 log,
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

func (uc *ModulesUseCase) QueueQuizletModuleImport(
	module *entity.Module,
	quizletModuleID string,
) error {
	importWork := &QuizletImportWork{
		repo:                uc.modulesRepo,
		quizletModuleParser: uc.quizletModuleParser,
		log:                 uc.log,
		quizletModuleID:     quizletModuleID,
		module:              module,
	}

	return uc.quizletImportWP.QueueWork(importWork)
}
