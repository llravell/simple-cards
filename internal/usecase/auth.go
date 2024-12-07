package usecase

import (
	"context"

	"github.com/llravell/simple-cards/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

const passwordCryptCost = 14

type AuthUseCase struct {
	repo      UserRepository
	jwtSecret []byte
}

func NewAuthUseCase(repo UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (auth *AuthUseCase) RegisterUser(ctx context.Context, login string, password string) (*entity.User, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCryptCost)
	if err != nil {
		return nil, err
	}

	return auth.repo.StoreUser(ctx, login, string(passwordBytes))
}

func (auth *AuthUseCase) BuildUserToken(user *entity.User) (string, error) {
	return entity.BuildJWTString(user.UUID, auth.jwtSecret)
}

func (auth *AuthUseCase) VerifyUser(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := auth.repo.FindUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
