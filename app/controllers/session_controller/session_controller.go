package session_controller

import (
	"Flashcards/app/controllers/common"
	"Flashcards/app/models"
	"Flashcards/app/services/session_service"
	"Flashcards/app/services/session_state_service"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type SessionController struct {
	service             *session_service.SessionService
	sessionStateService *session_state_service.SessionStateService
}

func New(service *session_service.SessionService, sessionStateService *session_state_service.SessionStateService) *SessionController {
	return &SessionController{
		service:             service,
		sessionStateService: sessionStateService,
	}
}

func (c *SessionController) Create(ctx *gin.Context) {
	var sessionDTO models.SessionDTO
	if err := ctx.ShouldBindJSON(&sessionDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := c.service.Create(sessionDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

func (c *SessionController) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "session.get.founded",
		NotFound:            "session.get.NotFound",
		BadRequest:          "session.get.BadRequest",
		InternalServerError: "session.get.Error",
	}

	id := ctx.Param("id")

	session, err := c.service.GetByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, messageTypes.NotFound, err))
			return
		}
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "FlashcardController",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: session,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SessionController) AnswerQuestion(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "session.get.founded",
		NotFound:            "session.get.NotFound",
		BadRequest:          "session.get.BadRequest",
		InternalServerError: "session.get.Error",
	}

	var answerDTO models.UserAnswerDTO
	if err := ctx.ShouldBindJSON(&answerDTO); err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	session, err := c.service.AnswerQuestion(answerDTO)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, messageTypes.NotFound, err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.MetaResponse{
		ObjectName: "FlashcardController",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: session,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

func (c *SessionController) GetState(context *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "session.get.founded",
		NotFound:            "session.get.NotFound",
		BadRequest:          "session.get.BadRequest",
		InternalServerError: "session.get.Error",
	}

	id := context.Param("id")

	state, err := c.sessionStateService.GetByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			common.SendResponse(context, http.StatusNotFound, models.KnownError(http.StatusNotFound, messageTypes.NotFound, err))
			return
		}
		common.SendResponse(context, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "FlashcardController",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: state,
	}

	common.SendResponse(context, http.StatusOK, response)
}
