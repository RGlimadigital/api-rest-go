package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model

	UserID    int   `json:"user_id"`
	PictureID []int `json:"picture_id"`
	Picture   []Picture
}

func newGallery(userID int, picturesID []int, picture []Picture) *Gallery {
	return &Gallery{
		UserID:    userID,
		PictureID: picturesID,
		Picture:   picture,
	}
}
