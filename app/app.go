package app

import (
	"context"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"mbsms-api/app/handler"
	"mbsms-api/app/service"
)

// New creates a new app instance
func New(sender service.Sender) *App {
	app := &App{SMSSender: sender}
	app.init()

	return app
}

// App has api router and sms worker instances
type App struct {
	SMSSender service.Sender

	api    *rest.Api
	server *http.Server
}

func (a *App) init() {
	a.api = rest.NewApi()
	a.api.Use(rest.DefaultDevStack...)
	a.setRouter()
}

func (a *App) setRouter() {
	router, err := rest.MakeRouter(
		rest.Post("/messages", a.PostMessage),
	)

	if err != nil {
		log.Fatal(err)
	}
	a.api.SetApp(router)
}

// PostMessage handles post message requests
func (a *App) PostMessage(w rest.ResponseWriter, r *rest.Request) {
	handler.PostMessage(a.SMSSender, w, r)
}

// Run starts serving the REST API
func (a *App) Run(host string) {
	log.Printf("Listening on http://0.0.0.0%s\n", host)
	a.server = &http.Server{Addr: host}
	a.server.Handler = a.api.MakeHandler()
	if err := a.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Shutdown gracefully shuts down the sms worker and
// API server without interrupting any active request
func (a *App) Shutdown(ctx context.Context) error {
	a.server.Shutdown(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
