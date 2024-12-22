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
	"github.com/llravell/simple-cards/pkg/quizlet"
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

	quizletParser, err := quizlet.NewParser()
	if err != nil {
		log.Fatalf("quizlet parser initialize error: %s", err)
	}

	logger := logger.Get()

	usersRepository := repository.NewUsersRepository(db)
	modulesRepository := repository.NewModulesRepository(db)
	cardsRepository := repository.NewCardsRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	healthUseCase := usecase.NewHealthUseCase(db)
	authUseCase := usecase.NewAuthUseCase(usersRepository, jwtManager)
	modulesUseCase := usecase.NewModulesUseCase(modulesRepository, cardsRepository, quizletParser)
	cardsUseCase := usecase.NewCardsUseCase(cardsRepository)

	app.New(
		healthUseCase,
		authUseCase,
		modulesUseCase,
		cardsUseCase,
		jwtManager,
		logger,
		app.Addr(cfg.Addr),
	).Run()
}
