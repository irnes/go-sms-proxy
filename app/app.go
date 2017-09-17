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
func New(sms *service.SMSService) *App {
	app := &App{sms: sms}
	app.init()

	return app
}

// App has api router and sms worker instances
type App struct {
	sms *service.SMSService

	api  *rest.Api
	http *http.Server
}

func (a *App) init() {
	a.api = rest.NewApi()
	a.api.Use(rest.DefaultProdStack...)
	a.setRouter()

	a.http = &http.Server{}
	a.http.Handler = a.api.MakeHandler()
}

func (a *App) setRouter() {
	router, err := rest.MakeRouter(
		rest.Post("/messages", handler.PostMessage(a.sms)),
	)

	if err != nil {
		log.Fatal(err)
	}
	a.api.SetApp(router)
}

// Run starts serving the REST API
func (a *App) Run(host string) {
	log.Printf("Listening on http://%s\n", host)
	a.http.Addr = host
	if err := a.http.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Shutdown gracefully shuts down the sms provider and
// API server without interrupting any active request
func (a *App) Shutdown(ctx context.Context) {
	//a.http.Shutdown(ctx) // requires golang >= v1.8

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
