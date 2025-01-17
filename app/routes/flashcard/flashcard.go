package flashcard

import (
	controller "Flashcards/app/controllers/flashcard_controller"
	service "Flashcards/app/services/flashcard_service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	flashCardService := service.New()
	flashCardController := controller.New(flashCardService)

	v1 := g.Group("/v1")
	{
		flashcards := v1.Group("/flashcards")
		{
			flashcards.GET("", flashCardController.Get)
			flashcards.POST("", flashCardController.Create)
			flashcards.GET("/:id", flashCardController.GetByID)
			flashcards.GET("/tag/:tag", flashCardController.GetByTag)
		}
	}
}
