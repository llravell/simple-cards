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
	healthUseCase := usecase.NewHealthUseCase(db)
	authUseCase := usecase.NewAuthUseCase(usersRepository, cfg.JWTSecret)

	app.New(
		healthUseCase,
		authUseCase,
		log,
		app.Addr(cfg.Addr),
		app.JWTSecret(cfg.JWTSecret),
	).Run()
}
