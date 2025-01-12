package usecase

import (
	"context"
	"errors"
	"reflect"
	"time"
)

const pingTimeout = 15 * time.Second

var ErrHasNotConnection = errors.New("has not db connection")

type HealthUseCase struct {
	repo HealthRepository
}

func NewHealthUseCase(repo HealthRepository) *HealthUseCase {
	return &HealthUseCase{repo: repo}
}

func (healthUC HealthUseCase) PingContext(ctx context.Context) error {
	v := reflect.ValueOf(healthUC.repo)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return ErrHasNotConnection
	}

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	return healthUC.repo.PingContext(ctx)
}
