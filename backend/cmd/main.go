package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MyHourly/TalentIQ/backend/database"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/handler"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/repository"
	"github.com/MyHourly/TalentIQ/backend/internal/skill-intelligence/service"
)

func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(
		w,
		"TalentIQ Backend is running",
	)
}

// responseWriter is used to capture the HTTP response status.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code returned by the API.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures responses where WriteHeader was not explicitly called.
func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	return rw.ResponseWriter.Write(data)
}

// requestLogger logs every API request in the terminal.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		fmt.Printf(
			"[API] %s %s -> %d %s (%v)\n",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			http.StatusText(rw.statusCode),
			duration,
		)
	})
}

func main() {

	// Create context with 10 second timeout.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Connect to PostgreSQL.
	db, err := database.NewConnection(ctx)

	if err != nil {
		fmt.Println(
			"Database connection failed:",
			err,
		)
		return
	}

	defer db.Close()

	fmt.Println(
		"Database connected successfully",
	)

	// --------------------------------------------------
	// Skill Repository
	// --------------------------------------------------

	skillRepository :=
		repository.NewPostgresSkillRepository(db)

	// --------------------------------------------------
	// Skill Service
	// --------------------------------------------------

	skillService :=
		service.NewSkillService(skillRepository)

	// --------------------------------------------------
	// Skill Handler
	// --------------------------------------------------

	skillHandler :=
		handler.NewSkillHandler(skillService)

	// --------------------------------------------------
	// Category Repository
	// --------------------------------------------------

	categoryRepository :=
		repository.NewPostgresCategoryRepository(db)

	// --------------------------------------------------
	// Category Service
	// --------------------------------------------------

	categoryService :=
		service.NewCategoryService(categoryRepository)

	// --------------------------------------------------
	// Category Handler
	// --------------------------------------------------

	categoryHandler :=
		handler.NewCategoryHandler(categoryService)

	// --------------------------------------------------
	// Health Check
	// --------------------------------------------------

	http.HandleFunc(
		"/health",
		healthHandler,
	)

	// --------------------------------------------------
	// Category Routes
	// --------------------------------------------------

	// POST and GET categories.
	http.HandleFunc(
		"/api/v1/categories",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodPost:
				categoryHandler.CreateCategory(w, r)

			case http.MethodGet:
				categoryHandler.GetAllCategories(w, r)

			default:
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	// --------------------------------------------------
	// Skill Routes
	// --------------------------------------------------

	// POST and GET all skills.
	http.HandleFunc(
		"/api/v1/skills",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodPost:
				skillHandler.CreateSkill(w, r)

			case http.MethodGet:
				skillHandler.GetAllSkills(w, r)

			default:
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	// --------------------------------------------------
	// Skill ID and Skill Category Routes
	// --------------------------------------------------

	http.HandleFunc(
		"/api/v1/skills/",
		func(w http.ResponseWriter, r *http.Request) {

			// Example:
			// /api/v1/skills/{skill_id}/categories

			if strings.HasSuffix(
				r.URL.Path,
				"/categories",
			) {

				switch r.Method {

				case http.MethodPost:
					categoryHandler.MapSkillToCategory(
						w,
						r,
					)

				case http.MethodGet:
					categoryHandler.GetCategoriesBySkillID(
						w,
						r,
					)

				default:
					http.Error(
						w,
						"Method not allowed",
						http.StatusMethodNotAllowed,
					)
				}

				return
			}

			// Example:
			// /api/v1/skills/{skill_id}

			switch r.Method {

			case http.MethodGet:
				skillHandler.GetSkillByID(w, r)

			case http.MethodPut:
				skillHandler.UpdateSkill(w, r)

			case http.MethodDelete:
				skillHandler.DeactivateSkill(w, r)

			default:
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	// --------------------------------------------------
	// Server
	// --------------------------------------------------

	fmt.Println(
		"TalentIQ Backend started",
	)

	fmt.Println(
		"Server running on http://localhost:8080",
	)

	fmt.Println(
		"Waiting for API requests...",
	)

	// Add request logging middleware.
	loggedHandler :=
		requestLogger(http.DefaultServeMux)

	err = http.ListenAndServe(
		":8080",
		loggedHandler,
	)

	if err != nil {
		fmt.Println(
			"Server error:",
			err,
		)
	}
}
