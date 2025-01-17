package models

type UserAnswer struct {
	StudentId   string `json:"studentId" bson:"studentId"`
	SessionId   string `json:"sessionId" bson:"sessionId"`
	CardId      string `json:"cardId" bson:"cardId"`
	AnswerIndex int    `json:"answerIndex" bson:"answerIndex"`
}

func (u UserAnswer) Collection() string {
	return "userAnswer"
}

type UserAnswerDTO struct {
	StudentId   string `json:"studentId" bson:"studentId"`
	SessionId   string `json:"sessionId" bson:"sessionId"`
	CardId      string `json:"cardId" bson:"cardId"`
	AnswerIndex int    `json:"answerIndex" bson:"answerIndex"`
}
