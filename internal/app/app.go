package app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/turbak/joom-calendar/internal/creating"
	"github.com/turbak/joom-calendar/internal/listing"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"github.com/turbak/joom-calendar/internal/pkg/http/mw"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
)

type Creator interface {
	CreateUser(ctx context.Context, user creating.User) (int, error)
	CreateEvent(ctx context.Context, event creating.Event) (int, error)
}

type Lister interface {
	GetEventByID(ctx context.Context, eventID int) (*listing.Event, error)
}

type Inviter interface {
	AcceptInvite(ctx context.Context, inviteID int) error
	DeclineInvite(ctx context.Context, inviteID int) error
}

type App struct {
	publicRouter chi.Router

	creator Creator
	lister  Lister
	inviter Inviter
}

func New(creator Creator, lister Lister, inviter Inviter) *App {
	a := &App{
		creator: creator,
		lister:  lister,
		inviter: inviter,
	}
	return a
}

func (a *App) Routes() chi.Router {
	a.publicRouter = chi.NewRouter()
	a.publicRouter.Use(mw.Recover(), mw.ResponseTimeLogging())

	a.publicRouter.Post("/users", httputil.Handler(a.handleCreateUser()))

	a.publicRouter.Post("/events", httputil.Handler(a.handleCreateEvent()))
	a.publicRouter.Get("/events/{event_id}", httputil.Handler(a.handleGetEvent()))

	a.publicRouter.Post("/event-invites/{invite_id}:accept", httputil.Handler(a.handleAcceptInvite()))
	a.publicRouter.Post("/event-invites/{invite_id}:decline", httputil.Handler(a.handleDeclineInvite()))

	return a.publicRouter
}

func (a *App) Run(addr string) error {
	logger.Debugf("app running on %s", addr)

	return http.ListenAndServe(addr, a.Routes())
}
