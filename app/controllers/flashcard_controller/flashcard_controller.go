package flashcard_controller

import (
	"Flashcards/app/controllers/common"
	"Flashcards/app/models"
	"Flashcards/app/services/flashcard_service"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type FlashcardController struct {
	FlashcardService *flashcard_service.FlashcardService
}

func New(flashcardService *flashcard_service.FlashcardService) *FlashcardController {
	return &FlashcardController{
		FlashcardService: flashcardService,
	}
}

func (f *FlashcardController) Get(ctx *gin.Context) {
	flashcards, err := f.FlashcardService.Get()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, flashcards)
}

func (f *FlashcardController) Create(ctx *gin.Context) {
	var createFlashcard models.FlashcardDTO
	messageTypes := &models.MessageTypes{
		Created:             "flashcard.Create.Created",
		BadRequest:          "flashcard.Create.BadRequest",
		InternalServerError: "flashcard.Create.Error",
	}

	if err := ctx.ShouldBindJSON(&createFlashcard); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	flashcard, err := f.FlashcardService.Create(createFlashcard)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "flashcard",
		TotalCount: 1,
		Offset:     0,
		Count:      1,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: flashcard,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

func (f *FlashcardController) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "flashcard.get.founded",
		NotFound:            "flashcard.get.NotFound",
		BadRequest:          "flashcard.get.BadRequest",
		InternalServerError: "flashcard.get.Error",
	}

	id := ctx.Param("id")

	flashcard, err := f.FlashcardService.GetById(id)
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
		Data: flashcard,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

func (f *FlashcardController) GetByTag(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "flashcard.get.founded",
		NotFound:            "flashcard.get.NotFound",
		BadRequest:          "flashcard.get.BadRequest",
		InternalServerError: "flashcard.get.Error",
	}

	tag := ctx.Param("tag")

	flashcards, err := f.FlashcardService.GetByTag(tag)
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
		TotalCount: len(flashcards),
		Count:      len(flashcards),
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: flashcards,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
