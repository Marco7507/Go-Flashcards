package main

import (
	controllers "Flashcards/app/controllers/common"
	routes "Flashcards/app/routes/common"

	"github.com/gin-gonic/gin"
)

// init the router
func setupRouter() *gin.Engine {
	router := routes.SetupRouter()
	router.GET("/ping", controllers.Ping)
	router.GET("/version", controllers.Version)

	return router
}
