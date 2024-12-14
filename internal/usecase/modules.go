package usecase

import (
	"context"

	"github.com/llravell/simple-cards/internal/entity"
)

type ModulesUseCase struct {
	repo ModuleRepository
}

func NewModulesUseCase(repo ModuleRepository) *ModulesUseCase {
	return &ModulesUseCase{
		repo: repo,
	}
}

func (uc *ModulesUseCase) GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error) {
	return uc.repo.GetAllModules(ctx, userUUID)
}

func (uc *ModulesUseCase) CreateNewModule(
	ctx context.Context,
	userUUID string,
	moduleName string,
) (*entity.Module, error) {
	return uc.repo.CreateNewModule(ctx, userUUID, moduleName)
}

func (uc *ModulesUseCase) UpdateModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
	moduleName string,
) (*entity.Module, error) {
	return uc.repo.UpdateModule(ctx, userUUID, moduleUUID, moduleName)
}

func (uc *ModulesUseCase) DeleteModule(
	ctx context.Context,
	userUUID string,
	moduleUUID string,
) error {
	return uc.repo.DeleteModule(ctx, userUUID, moduleUUID)
}
