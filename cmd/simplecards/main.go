package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/llravell/simple-cards/config"
	"github.com/llravell/simple-cards/internal/app"
	"github.com/llravell/simple-cards/internal/repository"
	"github.com/llravell/simple-cards/internal/usecase"
	"github.com/llravell/simple-cards/logger"
	"github.com/llravell/simple-cards/pkg/auth"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	var db *sql.DB

	if cfg.DatabaseURI != "" {
		db, err = sql.Open("pgx", cfg.DatabaseURI)
		if err != nil {
			log.Fatalf("open db error: %s", err)
		}
		defer db.Close()
	}

	log := logger.Get()

	usersRepository := repository.NewUsersRepository(db)
	modulesRepository := repository.NewModulesRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	healthUseCase := usecase.NewHealthUseCase(db)
	authUseCase := usecase.NewAuthUseCase(usersRepository, jwtManager)
	modulesUseCase := usecase.NewModulesUseCase(modulesRepository)

	app.New(
		healthUseCase,
		authUseCase,
		modulesUseCase,
		jwtManager,
		log,
		app.Addr(cfg.Addr),
	).Run()
}
