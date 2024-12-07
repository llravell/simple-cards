package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/llravell/simple-cards/internal/controller/http/health"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/internal/usecase"
	"github.com/rs/zerolog"
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
	healthUseCase *usecase.HealthUseCase
	router        chi.Router
	log           zerolog.Logger
	addr          string
	jwtSecret     string
}

func Addr(addr string) Option {
	return func(app *App) {
		app.addr = addr
	}
}

func JWTSecret(secret string) Option {
	return func(app *App) {
		app.jwtSecret = secret
	}
}

func New(
	healthUseCase *usecase.HealthUseCase,
	log zerolog.Logger,
	opts ...Option,
) *App {
	app := &App{
		healthUseCase: healthUseCase,
		log:           log,
		router:        chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *App) Run() {
	healthRoutes := health.NewHealthRoutes(app.healthUseCase, app.log)

	app.router.Use(middleware.LoggerMiddleware(app.log))
	healthRoutes.Apply(app.router)

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
