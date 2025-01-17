package session_service

import (
	"Flashcards/app/functions"
	"Flashcards/app/models"
	"Flashcards/app/server"
	"Flashcards/app/services/flashcard_service"
	"Flashcards/app/services/session_state_service"
	"Flashcards/app/services/student_service"
	"Flashcards/app/services/user_answer_service"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

const NumFlashcards = 5

type SessionService struct {
	validate            *validator.Validate
	StudentService      *student_service.StudentService
	FlashcardService    *flashcard_service.FlashcardService
	SessionStateService *session_state_service.SessionStateService
	UserAnswerService   *user_answer_service.UserAnswerService
}

func New() *SessionService {
	return &SessionService{
		validate:            validator.New(),
		StudentService:      student_service.New(),
		FlashcardService:    flashcard_service.New(),
		SessionStateService: session_state_service.New(),
		UserAnswerService:   user_answer_service.New(),
	}
}

func (s *SessionService) Create(sessionDTO models.SessionDTO) (*models.Session, error) {
	var session models.Session
	var err error

	srv := server.GetServer()
	collection := srv.Database.Collection(session.Collection())

	err = functions.ConvertInputStructToDataStruct(sessionDTO, &session)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	student, err := s.StudentService.GetByID(sessionDTO.StudentId)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	flashcards, err := s.FlashcardService.GetRandomsByTag(sessionDTO.Category, NumFlashcards)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	flashcardsIds := make([]string, 0)
	for _, flashcard := range flashcards {
		flashcardsIds = append(flashcardsIds, flashcard.ID)
	}

	session.ID = functions.NewUUID()
	session.StudentId = student.CustomID
	session.Flashcards = flashcardsIds

	_, err = collection.InsertOne(context.TODO(), session)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	var currentCardId string
	if len(flashcards) > 0 {
		currentCardId = flashcards[0].ID
	}

	sessionState := models.SessionState{
		SessionId:     session.ID,
		CurrentCardId: currentCardId,
	}

	_, err = s.SessionStateService.Create(sessionState)

	return &session, nil
}

func (s *SessionService) GetByID(id string) (*models.Session, error) {
	var session models.Session

	srv := server.GetServer()
	collection := srv.Database.Collection(session.Collection())

	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&session)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) AnswerQuestion(dto models.UserAnswerDTO) (sessionState *models.SessionState, err error) {
	session, err := s.GetByID(dto.SessionId)
	if err != nil {
		log.Error().Err(err).Msg("Session not found")
		return nil, err
	}

	student, err := s.StudentService.GetByID(session.StudentId)
	if err != nil {
		log.Error().Err(err).Msg("Student not found")
		return nil, err
	}

	sessionState, err = s.SessionStateService.GetByID(session.ID)
	if err != nil {
		log.Error().Err(err).Msg("Session state missing")
		return nil, err
	}
	if sessionState.IsFinished {
		log.Error().Err(err).Msg("Session is finished")
		return nil, err
	}

	flashcard, err := s.FlashcardService.GetById(sessionState.CurrentCardId)
	if err != nil {
		log.Error().Err(err).Msg("Flashcard not found")
		return nil, err
	}

	if dto.AnswerIndex < 0 || dto.AnswerIndex >= len(session.Flashcards) {
		log.Error().Err(err).Msg("Invalid answer index")
		return nil, errors.New("invalid answer index")
	}

	var currentCardId string
	if sessionState.CurrentCardIndex+1 < len(session.Flashcards) {
		currentCardId = session.Flashcards[sessionState.CurrentCardIndex+1]
	}

	sessionState.CurrentCardId = currentCardId
	sessionState.CurrentCardIndex = sessionState.CurrentCardIndex + 1
	sessionState.IsFinished = sessionState.CurrentCardIndex == len(session.Flashcards)-1

	if flashcard.CorrectAnswerIndex == dto.AnswerIndex {
		sessionState.Score++
	}

	_, err = s.SessionStateService.Update(sessionState)
	if err != nil {
		log.Error().Err(err).Msg("Error updating session state")
		return nil, err
	}

	userAnswer := models.UserAnswer{
		SessionId:   session.ID,
		StudentId:   student.CustomID,
		CardId:      currentCardId,
		AnswerIndex: dto.AnswerIndex,
	}

	_, err = s.UserAnswerService.Create(userAnswer)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return sessionState, nil
}
