package models

type Session struct {
	ID         string   `json:"id" bson:"id"`
	StudentId  string   `json:"studentId" bson:"studentId"`
	Category   string   `json:"category" bson:"category"`
	Flashcards []string `json:"flashcards" bson:"flashcards"`
}

func (s Session) Collection() string {
	return "sessions"
}

type SessionDTO struct {
	StudentId string `json:"studentId" binding:"required"`
	Category  string `json:"category" binding:"required"`
}
