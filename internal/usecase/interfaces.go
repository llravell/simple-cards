package usecase

import (
	"context"
	"time"

	"github.com/llravell/simple-cards/internal/entity"
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

	ModuleRepository interface {
		GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error)
		CreateNewModule(ctx context.Context, userUUID string, moduleName string) (*entity.Module, error)
		UpdateModule(ctx context.Context, userUUID string, moduleUUID string, moduleName string) (*entity.Module, error)
		DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error
	}

	JWTIssuer interface {
		Issue(userUUID string, ttl time.Duration) (string, error)
	}
)
