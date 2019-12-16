package models

import "github.com/jinzhu/gorm"

import "time"

type Comment struct {
	gorm.Model

	UserID    int       `json:"user_id"`
	PictureID int       `json:"picture_id"`
	Created   time.Time `json:"date"`
	Comment   string    `json: "comment"`
	User      User
}

func NewComment(userID int, pictureID int, created time.Time, comment string, user User) *Comment {
	return &Comment{
		UserID:    userID,
		PictureID: pictureID,
		Created:   created,
		Comment:   comment,
		User:      user,
	}
}
