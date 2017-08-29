package app

import (
	"context"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// New creates a new app instance
func New(sender Sender) *App {
	app := &App{SMSSender: sender}
	app.init()

	return app
}

// App has api router and sms worker instances
type App struct {
	SMSSender Sender

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
	payload := &BaseMessage{}
	err := r.DecodeJsonPayload(payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respChan := a.SMSSender.Send(payload)
	response := <-respChan

	w.WriteJson(&response)
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
