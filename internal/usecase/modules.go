package usecase

import (
	"context"
	"io"

	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

type ModulesUseCase struct {
	modulesRepo         ModulesRepository
	cardsRepo           CardsRepository
	quizletModuleParser QuizletModuleParser
	quizletImportWP     QuizletImportWorkerPool
	csvImportWP         CSVImportWorkerPool
	log                 *zerolog.Logger
}

func NewModulesUseCase(
	modulesRepo ModulesRepository,
	cardsRepo CardsRepository,
	quizletModuleParser QuizletModuleParser,
	quizletImportWP QuizletImportWorkerPool,
	csvImportWP CSVImportWorkerPool,
	log *zerolog.Logger,
) *ModulesUseCase {
	return &ModulesUseCase{
		modulesRepo:         modulesRepo,
		cardsRepo:           cardsRepo,
		quizletModuleParser: quizletModuleParser,
		quizletImportWP:     quizletImportWP,
		csvImportWP:         csvImportWP,
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

func (uc *ModulesUseCase) QueueCSVModuleImport(
	module *entity.Module,
	reader io.ReadCloser,
) error {
	importWork := &CSVImportWork{
		repo:   uc.modulesRepo,
		log:    uc.log,
		module: module,
		reader: reader,
	}

	return uc.csvImportWP.QueueWork(importWork)
}
