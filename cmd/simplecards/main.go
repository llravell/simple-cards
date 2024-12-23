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
	"github.com/llravell/simple-cards/pkg/workerpool"
)

const quizletImportWorkersAmount = 4

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	quizletParser, err := quizlet.NewParser()
	if err != nil {
		log.Fatalf("quizlet parser initialize error: %s", err)
	}

	var db *sql.DB

	if cfg.DatabaseURI != "" {
		db, err = sql.Open("pgx", cfg.DatabaseURI)
		if err != nil {
			log.Fatalf("open db error: %s", err)
		}
		defer db.Close()
	}

	logger := logger.Get()

	usersRepository := repository.NewUsersRepository(db)
	modulesRepository := repository.NewModulesRepository(db)
	cardsRepository := repository.NewCardsRepository(db)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	quizletImportWorkerPool := workerpool.New[*usecase.QuizletImportWork](quizletImportWorkersAmount)

	healthUseCase := usecase.NewHealthUseCase(db)
	authUseCase := usecase.NewAuthUseCase(usersRepository, jwtManager)
	modulesUseCase := usecase.NewModulesUseCase(
		modulesRepository,
		cardsRepository,
		quizletParser,
		quizletImportWorkerPool,
		&logger,
	)
	cardsUseCase := usecase.NewCardsUseCase(cardsRepository)

	quizletImportWorkerPool.ProcessQueue()

	defer func() {
		quizletImportWorkerPool.Close()

		logger.Info().Msg("quizlet import worker pool closing...")
		quizletImportWorkerPool.Wait()
	}()

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
