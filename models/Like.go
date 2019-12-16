package models

import "time"

type Likes struct {
	UserID    int       `json:"user_id"`
	PictureID int       `json:"picture_id"`
	Liked     time.Time `json:"date"`
}

func NewLike(userID int, pictureID int, liked time.Time) *Likes {
	return &Likes{
		UserID:    userID,
		PictureID: pictureID,
		Liked:     liked,
	}
}
