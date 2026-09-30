package main

import (
	"context"
	"fmt"
	"net/http"
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

	// Create Skill Repository.
	skillRepository :=
		repository.NewPostgresSkillRepository(db)

	// Create Skill Service.
	skillService :=
		service.NewSkillService(skillRepository)

	// Create Skill Handler.
	skillHandler :=
		handler.NewSkillHandler(skillService)

	// Health check.
	http.HandleFunc(
		"/health",
		healthHandler,
	)

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

	// GET, PUT and DELETE one skill.
	http.HandleFunc(
		"/api/v1/skills/",
		func(w http.ResponseWriter, r *http.Request) {

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

	fmt.Println(
		"TalentIQ Backend started",
	)

	fmt.Println(
		"Server running on http://localhost:8080",
	)

	err = http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		fmt.Println(
			"Server error:",
			err,
		)
	}
}
