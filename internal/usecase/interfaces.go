package usecase

import (
	"context"
	"time"

	"github.com/llravell/simple-cards/internal/entity"
	"github.com/llravell/simple-cards/pkg/quizlet"
)

//go:generate ../../bin/mockgen -source=interfaces.go -destination=../mocks/mock_usecase.go -package=mocks

type (
	HealthRepository interface {
		PingContext(ctx context.Context) error
	}

	UserRepository interface {
		StoreUser(ctx context.Context, login string, password string) (*entity.User, error)
		FindUserByLogin(ctx context.Context, login string) (*entity.User, error)
	}

	ModulesRepository interface {
		GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error)
		GetModule(ctx context.Context, userUUID string, moduleUUID string) (*entity.Module, error)
		CreateNewModule(ctx context.Context, userUUID string, moduleName string) (*entity.Module, error)
		CreateNewModuleWithCards(ctx context.Context, moduleWithCards *entity.ModuleWithCards) error
		UpdateModule(ctx context.Context, userUUID string, moduleUUID string, moduleName string) (*entity.Module, error)
		DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error
		ModuleExists(ctx context.Context, userUUID string, moduleUUID string) (bool, error)
	}

	CardsRepository interface {
		GetModuleCards(ctx context.Context, moduleUUID string) ([]*entity.Card, error)
		CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
		SaveCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
		DeleteCard(ctx context.Context, moduleUUID string, cardUUID string) error
	}

	JWTIssuer interface {
		Issue(userUUID string, ttl time.Duration) (string, error)
	}

	QuizletModuleParser interface {
		Parse(ctx context.Context, moduleID string) ([]quizlet.Card, error)
	}

	QuizletImportWorkerPool interface {
		QueueWork(w *QuizletImportWork) error
	}

	CSVImportWorkerPool interface {
		QueueWork(w *CSVImportWork) error
	}
)
