package http

import (
	"context"
	"io"
	"time"

	"github.com/llravell/simple-cards/internal/entity"
)

type HealthUseCase interface {
	PingContext(ctx context.Context) error
}

type AuthUseCase interface {
	VerifyUser(ctx context.Context, login string, password string) (*entity.User, error)
	BuildUserToken(user *entity.User, ttl time.Duration) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (*entity.User, error)
}

type ModulesUseCase interface {
	GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error)
	GetModuleWithCards(ctx context.Context, userUUID string, moduleUUID string) (*entity.ModuleWithCards, error)
	CreateNewModule(ctx context.Context, userUUID string, moduleName string) (*entity.Module, error)
	UpdateModule(ctx context.Context, userUUID string, moduleUUID string, moduleName string) (*entity.Module, error)
	DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error
	ModuleExists(ctx context.Context, userUUID string, moduleUUID string) (bool, error)
	QueueQuizletModuleImport(module *entity.Module, quizletModuleID string) error
	QueueCSVModuleImport(module *entity.Module, reader io.ReadCloser) error
}

type CardsUseCase interface {
	GetModuleCards(ctx context.Context, moduleUUID string) ([]*entity.Card, error)
	CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
	SaveCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
	DeleteCard(ctx context.Context, moduleUUID string, cardUUID string) error
}
