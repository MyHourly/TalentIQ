package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	return router
}
