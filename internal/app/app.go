package app

import (
	"github.com/go-chi/chi"
	"github.com/turbak/joom-calendar/internal/pkg/http/mw"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
)

type App struct {
	publicRouter chi.Router
}

func New() *App {
	a := &App{
		publicRouter: chi.NewRouter(),
	}
	return a
}

func (a *App) Routes() chi.Router {
	a.publicRouter.Use(mw.Recover(), mw.ResponseTimeLogging(), mw.RequestBodyLogging())
	return a.publicRouter
}

func (a *App) Run(addr string) error {
	logger.Debugf("app running on %s", addr)

	return http.ListenAndServe(addr, a.Routes())
}
