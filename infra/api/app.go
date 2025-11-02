package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"

	echo "github.com/labstack/echo/v4"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
)

type App struct {
	router       *echo.Echo
	config       *config.Application
	logger       logger.Logger
	repositories repository.Collection
}

func NewApp(
	logger logger.Logger,
	config *config.Application,
	repositories repository.Collection,
) *App {
	logger.Info("initializing the HTTP server...")

	app := &App{
		router:       echo.New(),
		config:       config,
		logger:       logger,
		repositories: repositories,
	}

	app.configure()

	return app
}

func (app *App) configure() {
	app.configureRoutes()
}

func (app *App) Start(wg *sync.WaitGroup) error {
	errChannel := make(chan error)

	go func(errChannel chan error) {
		address := fmt.Sprintf("%s:%d", app.config.API.Host, app.config.API.Port)
		app.logger.Info("starting the  HTTP server...")
		err := app.router.Start(address)
		if err != nil {
			errChannel <- fmt.Errorf("error while starting the HTTP server: %w", err)
		}
	}(errChannel)

	go func() {
		for app.router.Listener == nil {
			time.Sleep(1 * time.Millisecond)
		}

		wg.Done()
	}()

	err := <-errChannel

	return err
}

func (app *App) Stop() {
	app.logger.Info("the HTTP server was stopped")
}
