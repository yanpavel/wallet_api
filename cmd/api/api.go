package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yanpavel/wallet_api/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr   string
	db     dbConfig
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Post("/", app.loginHandler)
		})
		r.Route("/register", func(r chi.Router) {
			r.Post("/", app.handleRegister)
		})
		r.Route("/wallet", func(r chi.Router) {
			r.Post("/", WithAuthJWT(app.changeBalanceHandler, app.store))

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", WithAuthJWT(app.getBalanceHandler, app.store))
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 1,
	}

	log.Printf("server has started at %v", app.config.addr)

	return srv.ListenAndServe()
}
