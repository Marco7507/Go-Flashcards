package session_state_service

import (
	"Flashcards/app/models"
	"Flashcards/app/server"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type SessionStateService struct {
}

func New() *SessionStateService {
	return &SessionStateService{}
}

func (s *SessionStateService) Create(sessionState models.SessionState) (*models.SessionState, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection(sessionState.Collection())

	_, err := collection.InsertOne(context.TODO(), sessionState)
	if err != nil {
		return nil, err
	}

	return &sessionState, nil
}

func (s *SessionStateService) GetByID(sessionId string) (*models.SessionState, error) {
	var sessionState models.SessionState
	srv := server.GetServer()
	collection := srv.Database.Collection(sessionState.Collection())

	err := collection.FindOne(context.TODO(), bson.D{{"sessionId", sessionId}}).Decode(&sessionState)
	if err != nil {
		return nil, err
	}

	return &sessionState, nil
}

func (s *SessionStateService) Update(sessionState *models.SessionState) (*models.SessionState, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection(sessionState.Collection())

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"sessionId": sessionState.SessionId},
		bson.D{
			{"$set", bson.D{
				{"currentCardId", sessionState.CurrentCardId},
				{"currentCardIndex", sessionState.CurrentCardIndex},
				{"isFinished", sessionState.IsFinished},
				{"score", sessionState.Score},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return sessionState, nil
}
