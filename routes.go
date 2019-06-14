package main

import (
	"github.com/alabianca/rapi-api/app"
	"github.com/alabianca/rapi-api/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func apiRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,          // Log API Requests
		app.JwtAuthentication,      // Check for presence of jwt token
		setupCORS().Handler,        // Allow Cross-Origin-Requests
		middleware.DefaultCompress, // Compress results
		middleware.RedirectSlashes, // Redirect slashes to no slash url versions
		middleware.Recoverer,       // recover from panic without crashing
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/user", userRoutes())
		r.Mount("/api/token", tokenRoutes())
		r.Mount("/api/resume", resumeRoutes())
	})

	return router
}

func setupCORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "UPDATE", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           500,
	})

	return c
}

func userRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", controllers.CreateUser)
	router.Post("/{userID}", controllers.PostUser)
	router.Get("/{userID}", controllers.GetUser)
	router.Get("/{userID}/records", controllers.GetRecordsForUser)

	return router
}

func tokenRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", controllers.CreateToken)

	return router
}

func resumeRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", controllers.CreateResume)

	return router
}

func recordRoutes() *chi.Mux {
	router := chi.NewRouter()

	return router
}
