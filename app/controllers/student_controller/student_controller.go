package student_controller

import (
	"Flashcards/app/controllers/common"
	"Flashcards/app/models"
	"Flashcards/app/services/student_service"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	StudentService *student_service.StudentService
}

func New(studentService *student_service.StudentService) *StudentController {
	return &StudentController{
		StudentService: studentService,
	}
}

// Get controller to get list of student
func (s *StudentController) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "student.Search.Found",
		BadRequest:          "student.Search.BadRequest",
		NotFound:            "student.Search.NotFound",
		InternalServerError: "student.Search.Error",
	}

	students, err := s.StudentService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
	}
	totalCount := len(students)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Data not found. ")))
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Available CountMax calculation
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Offset cannot be higher than count. ")))
	}

	sendingStudents := students[low:high]

	meta := models.MetaResponse{
		ObjectName: "StudentService",
		TotalCount: totalCount,
		Count:      len(sendingStudents),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingStudents,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Create controller to create new student
func (s *StudentController) Create(ctx *gin.Context) {
	var in models.StudentInput
	messageTypes := &models.MessageTypes{
		Created:             "student.Create.Created",
		BadRequest:          "student.Create.BadRequest",
		InternalServerError: "student.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	student, err := s.StudentService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "StudentService",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: student,
	}

	common.SendResponse(ctx, http.StatusCreated, response)
}

// GetByID controller to get one student by id
func (s *StudentController) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "student.get.founded",
		NotFound:            "student.get.NotFound",
		BadRequest:          "student.get.BadRequest",
		InternalServerError: "student.get.Error",
	}

	id := ctx.Param("id")

	student, err := s.StudentService.GetByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, messageTypes.NotFound, err))
			return
		}
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	// Créer la réponse avec le client
	meta := models.MetaResponse{
		ObjectName: "StudentService",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: student,
	}

	// Envoyer la réponse
	common.SendResponse(ctx, http.StatusOK, response)
}

// Update controller to update student
func (s *StudentController) Update(ctx *gin.Context) {
	var in models.StudentInput
	messageTypes := &models.MessageTypes{
		OK:                  "student.Update.Updated",
		BadRequest:          "student.Update.BadRequest",
		InternalServerError: "students.Update.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	id := ctx.Param("id")
	err := s.StudentService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "student updated"))
}

func (s *StudentController) Suspend(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "student.Suspend.Updated",
		InternalServerError: "student.Suspend.Error",
	}
	id := ctx.Param("id")
	err := s.StudentService.Suspend(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "student suspended"))
}

// GetByIDs To get all student customID
func (s *StudentController) GetByIDs(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "student.Search.Updated",
		NotFound:            "student.Search.NotFound",
		BadRequest:          "student.Search.BadRequest",
		InternalServerError: "student.Search.Error",
	}

	// Extraction des identifiants de la requête URL
	ids := ctx.Param("ids")
	idList := strings.Split(ids, "&")

	students, err := s.StudentService.GetByIds(idList)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	totalCount := len(students)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Data not found.")))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Calcul du CountMax disponible
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Offset cannot be higher than count.")))
		return
	}

	sendingStudent := students[low:high]

	meta := models.MetaResponse{
		ObjectName: "StudentService",
		TotalCount: totalCount,
		Count:      len(sendingStudent),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingStudent,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
