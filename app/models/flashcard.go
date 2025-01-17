package models

import "time"

type Flashcard struct {
	ID                 string    `json:"id" bson:"id"`
	Question           string    `json:"question" bson:"question"`
	Answers            []string  `json:"answers" bson:"answers"`
	CorrectAnswerIndex int       `json:"correctAnswerIndex" bson:"correctAnswerIndex"`
	Tags               []string  `json:"tags" bson:"tags"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
}

func (f *Flashcard) Collection() string {
	return "flashcard"
}

type FlashcardDTO struct {
	Question           string   `json:"question" bson:"question"`
	Answers            []string `json:"answers" bson:"answers"`
	CorrectAnswerIndex int      `json:"correctAnswerIndex" bson:"correctAnswerIndex"`
	Tags               []string `json:"tags" bson:"tags"`
}
