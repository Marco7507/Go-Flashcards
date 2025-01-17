package models

import "time"

// Student represent a student
type Student struct {
	CustomID    string    `bson:"customID" json:"customId"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	GSMNumber   string    `json:"gsmNumber"`
	Newsletter  bool      `json:"newsletter"`
	CreatedAt   time.Time `json:"creationDate"`
	Suspended   bool      `json:"suspended"`
}

// Collection Mongodb collection
func (s *Student) Collection() string {
	return "student"
}

// StudentInput
type StudentInput struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	GsmNumber   string `json:"gsmNumber"`
	Newsletter  bool   `json:"newsletter"`
}
