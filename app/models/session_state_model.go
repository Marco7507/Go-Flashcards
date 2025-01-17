package models

type SessionState struct {
	SessionId        string `json:"sessionId" bson:"sessionId"`
	CurrentCardIndex int    `json:"currentCardIndex" bson:"currentCardIndex"`
	CurrentCardId    string `json:"currentCardId" bson:"currentCardId"`
	IsFinished       bool   `json:"isFinished" bson:"isFinished"`
	Score            int    `json:"score" bson:"score"`
}

func (s SessionState) Collection() string {
	return "sessionState"
}
