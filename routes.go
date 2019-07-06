package main

import (
	"github.com/alabianca/rapi-api/app"
	"github.com/alabianca/rapi-api/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func apiRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		setupCORS().Handler,        // Allow Cross-Origin-Requests
		middleware.Logger,          // Log API Requests
		middleware.DefaultCompress, // Compress results
		middleware.RedirectSlashes, // Redirect slashes to no slash url versions
		middleware.Recoverer,       // recover from panic without crashing
	)

	router.Route("/v1", func(r chi.Router) {
		r.Use(app.JwtAuthentication)
		r.Mount("/api/user", userRoutes(api))
		r.Mount("/api/token", tokenRoutes(api))
		r.Mount("/api/resume", resumeRoutes(api))
		r.Mount("/api/key", keyRoutes(api))
		r.Mount("/api/metrics", metricRoutes(api))
	})

	router.Route("/pub/v1", func(r chi.Router) {
		r.Mount("/record", recordRoutes(api))
	})

	return router
}

func setupCORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "UPDATE", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           500,
	})

	return c
}

func userRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", api.CreateUser)
	router.Post("/{userID}", api.PostUser)
	router.Get("/{userID}", api.GetUser)
	router.Get("/{userID}/records", api.GetRecordsForUser)

	return router
}

func tokenRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", api.CreateToken)

	return router
}

func resumeRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", api.CreateResume)
	router.Get("/", api.GetResumes)

	return router
}

func keyRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/{resumeID}", api.CreateKey)
	router.Get("/{resumeID}", api.GetKeys)
	router.Patch("/{keyID}", api.PatchKey)
	router.Delete("/{keyID}", api.DeleteKey)

	return router
}

func metricRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", api.GetMetrics)

	return router
}

func recordRoutes(api *controllers.API) *chi.Mux {
	router := chi.NewRouter()

	r := router.With(app.CheckKey, app.LogKey(api))
	r.Get("/{resumeID}", api.GetResume)
	r.Get("/{resumeID}/experience", api.GetExperience)
	r.Get("/{resumeID}/education", api.GetEducation)
	r.Get("/{resumeID}/personal", api.GetPersonal)
	r.Get("/{resumeID}/projects", api.GetProjects)
	r.Get("/{resumeID}/skills", api.GetSkills)

	r.Post("/{resumeID}", api.GetResume)

	return router
}
