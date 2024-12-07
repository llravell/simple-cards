package usecase

import "context"

//go:generate ../../bin/mockgen -source=interfaces.go -destination=../mocks/mock_usecase.go -package=mocks

type (
	HealthRepository interface {
		PingContext(ctx context.Context) error
	}
)
