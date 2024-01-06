package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	if a.debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(a.LoadSession)

	//register routes
	mux.Get("/", a.homeHandler)
	mux.Get("/comments/{postId}", a.commentHandler)
	mux.Post("/comments/{postId}", a.commentPostHandler)

	mux.Get("/login", a.loginHandler)
	mux.Post("/login", a.loginPostHandler)
	mux.Get("/signup", a.signUpHandler)
	mux.Post("/signup", a.signPostUpHandler)
	mux.Get("/logout", a.authRequired(a.logoutHandler))

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
