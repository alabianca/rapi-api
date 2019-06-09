package main

import (
	"github.com/alabianca/rapi-api/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func apiRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,          // Log API Requests
		middleware.DefaultCompress, // Compress results
		middleware.RedirectSlashes, // Redirect slashes to no slash url versions
		middleware.Recoverer,       // recover from panic without crashing
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/user", userRoutes())
		r.Mount("/api/token", tokenRoutes())
	})

	return router
}

func userRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", controllers.CreateUser)
	router.Get("/{userID}", controllers.GetUser)

	return router
}

func tokenRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", controllers.CreateToken)

	return router
}
