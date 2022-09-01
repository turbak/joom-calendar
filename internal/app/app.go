package app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/turbak/joom-calendar/internal/adding"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"github.com/turbak/joom-calendar/internal/pkg/http/mw"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
)

type AddingService interface {
	CreateUser(ctx context.Context, user adding.User) (int, error)
	CreateEvent(ctx context.Context, event adding.Event) (int, error)
}

type App struct {
	publicRouter chi.Router

	addingService AddingService
}

func New(addingService AddingService) *App {
	a := &App{
		publicRouter:  chi.NewRouter(),
		addingService: addingService,
	}
	return a
}

func (a *App) Routes() chi.Router {
	a.publicRouter.Use(mw.Recover(), mw.ResponseTimeLogging())

	a.publicRouter.Post("/users", httputil.Handler(a.handleCreateUser()))

	a.publicRouter.Post("/events", httputil.Handler(a.handleCreateEvent()))
	return a.publicRouter
}

func (a *App) Run(addr string) error {
	logger.Debugf("app running on %s", addr)

	return http.ListenAndServe(addr, a.Routes())
}
