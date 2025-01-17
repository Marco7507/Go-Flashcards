package session

import (
	"Flashcards/app/controllers/session_controller"
	"Flashcards/app/services/session_service"
	"Flashcards/app/services/session_state_service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	sessionController := session_controller.New(
		session_service.New(),
		session_state_service.New())

	v1 := g.Group("/v1")
	{
		sessions := v1.Group("/sessions")
		{
			sessions.POST("", sessionController.Create)
			sessions.GET("/:id", sessionController.GetByID)
			sessions.POST("/:id/answer", sessionController.AnswerQuestion)
			sessions.GET("/:id/state", sessionController.GetState)
		}
	}
}
