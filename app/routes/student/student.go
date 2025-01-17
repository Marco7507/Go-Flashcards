package student

import (
	controller "Flashcards/app/controllers/student_controller"
	service "Flashcards/app/services/student_service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	servicesStudent := service.New()
	studentController := controller.New(servicesStudent)

	v1 := g.Group("/v1")
	{
		students := v1.Group("/students")
		{
			students.POST("", studentController.Create)
			students.GET("", studentController.Get)
			students.GET("/:id", studentController.GetByID)
			students.POST("/:id", studentController.Update)
			students.POST("/:id/suspend", studentController.Suspend)
			students.GET("/IDS/:ids", studentController.GetByIDs)
		}
	}
}
