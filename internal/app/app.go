package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/llravell/simple-cards/docs" //nolint:blank-imports
	httpCommon "github.com/llravell/simple-cards/internal/controller/http"
	"github.com/llravell/simple-cards/internal/controller/http/auth"
	"github.com/llravell/simple-cards/internal/controller/http/cards"
	"github.com/llravell/simple-cards/internal/controller/http/health"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/internal/controller/http/modules"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func startServer(addr string, handler http.Handler) error {
	server := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	return server.ListenAndServe()
}

type Option func(app *App)

type App struct {
	healthUseCase  httpCommon.HealthUseCase
	authUseCase    httpCommon.AuthUseCase
	modulesUseCase httpCommon.ModulesUseCase
	cardsUseCase   httpCommon.CardsUseCase
	jwtParser      middleware.JWTParser
	router         chi.Router
	log            zerolog.Logger
	addr           string
}

func Addr(addr string) Option {
	return func(app *App) {
		app.addr = addr
	}
}

func New(
	healthUseCase httpCommon.HealthUseCase,
	authUseCase httpCommon.AuthUseCase,
	modulesUseCase httpCommon.ModulesUseCase,
	cardsUseCase httpCommon.CardsUseCase,
	jwtParser middleware.JWTParser,
	log zerolog.Logger,
	opts ...Option,
) *App {
	app := &App{
		healthUseCase:  healthUseCase,
		authUseCase:    authUseCase,
		modulesUseCase: modulesUseCase,
		cardsUseCase:   cardsUseCase,
		jwtParser:      jwtParser,
		log:            log,
		router:         chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(app)
	}

	app.setupRoutes()

	return app
}

func (app *App) setupRoutes() {
	healthRoutes := health.NewRoutes(app.healthUseCase, app.log)
	authRoutes := auth.NewRoutes(app.authUseCase, app.log)
	modulesRoutes := modules.NewRoutes(app.modulesUseCase, app.log)
	cardsRoutes := cards.NewRoutes(app.modulesUseCase, app.cardsUseCase, app.log)

	app.router.Use(middleware.LoggerMiddleware(app.log))
	healthRoutes.Apply(app.router)
	authRoutes.Apply(app.router)

	app.router.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(app.jwtParser, app.log))

		modulesRoutes.Apply(r)
		cardsRoutes.Apply(r)
	})

	app.router.Get("/swagger/*", httpSwagger.Handler())
}

// Swagger spec:
// @title       Simple Cards API
// @version     1.0
// @host        localhost:8080
// @BasePath    /

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization
func (app *App) Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	serverNotify := make(chan error, 1)
	go func() {
		serverNotify <- startServer(app.addr, app.router)
		close(serverNotify)
	}()

	app.log.Info().
		Str("addr", app.addr).
		Msgf("starting server on '%s'", app.addr)

	select {
	case s := <-interrupt:
		app.log.Info().Str("signal", s.String()).Msg("interrupt")
	case err := <-serverNotify:
		app.log.Error().Err(err).Msg("server has been closed")
	}
}
