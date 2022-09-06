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
	"time"
)

type Creator interface {
	CreateUser(ctx context.Context, user creating.User) (int, error)
	CreateEvent(ctx context.Context, event creating.Event) (int, error)
}

type Lister interface {
	GetEventByID(ctx context.Context, eventID int) (*listing.Event, error)
	ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]listing.Event, error)
	GetNearestEmptyTimeInterval(ctx context.Context, userIDs []int, minDuration time.Duration) (time.Time, time.Time, error)
}

type Inviter interface {
	AcceptInvite(ctx context.Context, inviteID int) error
	DeclineInvite(ctx context.Context, inviteID int) error
}

type Authenticator interface {
	AuthenticateGithub(ctx context.Context, code string) (string, error)
	Middleware() func(http.Handler) http.Handler
}

type App struct {
	publicRouter chi.Router

	creator       Creator
	lister        Lister
	inviter       Inviter
	authenticator Authenticator
}

func New(creator Creator, lister Lister, inviter Inviter, authenticator Authenticator) *App {
	a := &App{
		creator:       creator,
		lister:        lister,
		inviter:       inviter,
		authenticator: authenticator,
	}
	return a
}

func (a *App) Routes() chi.Router {
	a.publicRouter = chi.NewRouter()
	a.publicRouter.Use(mw.Recover(), mw.ResponseTimeLogging())

	a.publicRouter.Group(func(r chi.Router) {
		r.Use(a.authenticator.Middleware())

		r.Post("/users", httputil.Handler(a.handleCreateUser()))

		r.Post("/events", httputil.Handler(a.handleCreateEvent()))
		r.Get("/events/{event_id}", httputil.Handler(a.handleGetEvent()))
		r.Get("/events:nearest-empty-time-interval", httputil.Handler(a.handleFindNearestTimeInterval()))

		r.Post("/event-invites/{invite_id}:accept", httputil.Handler(a.handleAcceptInvite()))
		r.Post("/event-invites/{invite_id}:decline", httputil.Handler(a.handleDeclineInvite()))

		r.Get("/users/{user_id}/events", httputil.Handler(a.handleGetUserEvents()))
	})

	a.publicRouter.Get("/login/github", a.handleGithubLogin())
	a.publicRouter.Get("/callbacks/github", httputil.Handler(a.handleGithubCallback()))

	return a.publicRouter
}

func (a *App) Run(addr string) error {
	logger.Debugf("app running on %s", addr)

	return http.ListenAndServe(addr, a.Routes())
}
