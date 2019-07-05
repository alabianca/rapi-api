package main

import (
	"github.com/alabianca/rapi-api/app"
	"github.com/alabianca/rapi-api/controllers"
	"github.com/alabianca/rapi-api/pub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func apiRoutes() *chi.Mux {
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
		r.Mount("/api/user", userRoutes())
		r.Mount("/api/token", tokenRoutes())
		r.Mount("/api/resume", resumeRoutes())
		r.Mount("/api/key", keyRoutes())
		r.Mount("/api/metrics", metricRoutes())
	})

	router.Route("/pub/v1", func(r chi.Router) {
		r.Mount("/record", recordRoutes())
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
	router.Get("/", controllers.GetResumes)

	return router
}

func keyRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/{resumeID}", controllers.CreateKey)
	router.Get("/{resumeID}", controllers.GetKeys)
	router.Patch("/{keyID}", controllers.PatchKey)
	router.Delete("/{keyID}", controllers.DeleteKey)

	return router
}

func metricRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", controllers.GetMetrics)

	return router
}

func recordRoutes() *chi.Mux {
	router := chi.NewRouter()

	r := router.With(app.CheckKey, app.LogKey)
	r.Get("/{resumeID}", pub.GetResume)
	r.Get("/{resumeID}/experience", pub.GetExperience)
	r.Get("/{resumeID}/education", pub.GetEducation)
	r.Get("/{resumeID}/personal", pub.GetPersonal)
	r.Get("/{resumeID}/projects", pub.GetProjects)
	r.Get("/{resumeID}/skills", pub.GetSkills)

	r.Post("/{resumeID}", pub.GetResume)

	return router
}
