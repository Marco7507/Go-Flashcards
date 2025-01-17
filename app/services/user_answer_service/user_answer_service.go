package user_answer_service

import (
	"Flashcards/app/models"
	"Flashcards/app/server"
	"context"
	"github.com/go-playground/validator/v10"
)

type UserAnswerService struct {
	validate *validator.Validate
}

func New() *UserAnswerService {
	return &UserAnswerService{
		validate: validator.New(),
	}
}

func (u *UserAnswerService) Create(userAnswer models.UserAnswer) (*models.UserAnswer, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection(userAnswer.Collection())

	_, err := collection.InsertOne(context.TODO(), userAnswer)
	if err != nil {
		return nil, err
	}

	return &userAnswer, nil
}

func (u *UserAnswerService) GetByID(sessionId, studentId, cardI string) (*models.UserAnswer, error) {
	var userAnswer models.UserAnswer
	srv := server.GetServer()
	collection := srv.Database.Collection(userAnswer.Collection())

	err := collection.FindOne(context.TODO(), models.UserAnswer{
		SessionId: sessionId,
		StudentId: studentId,
		CardId:    cardI,
	}).Decode(&userAnswer)
	if err != nil {
		return nil, err
	}

	return &userAnswer, nil
}

func (u *UserAnswerService) Update(userAnswer *models.UserAnswer) (*models.UserAnswer, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection(userAnswer.Collection())

	_, err := collection.UpdateOne(context.TODO(), models.UserAnswer{
		SessionId: userAnswer.SessionId,
		StudentId: userAnswer.StudentId,
		CardId:    userAnswer.CardId,
	}, userAnswer)
	if err != nil {
		return nil, err
	}

	return userAnswer, nil
}
