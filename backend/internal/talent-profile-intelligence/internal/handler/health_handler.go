package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Health handles the basic application health check.
//
// This endpoint only checks whether the application process
// is running. It does not depend on PostgreSQL.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "talent-profile-intelligence",
	})
}

// Ready handles the application readiness check.
//
// This endpoint checks whether the application can communicate
// with PostgreSQL.
func Ready(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Give the database a short amount of time to respond.
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		// Check the PostgreSQL connection.
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "NOT_READY",
				"service": "talent-profile-intelligence",
				"database": "DOWN",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "READY",
			"service": "talent-profile-intelligence",
			"database": "UP",
		})
	}
}