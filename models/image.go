package models

import "github.com/jinzhu/gorm"

type Image struct {
	gorm.Model
	UserID  int    `json:"user_id"`
	Thumb   string `json:"thumb_url"`
	Lowres  string `json:"lowres_url"`
	Highres string `json:"highres_url"`
}

func NewImage(userID int, thumb string, lowres string, highres string) *Image {
	return &Image{
		UserID:  userID,
		Thumb:   thumb,
		Lowres:  lowres,
		Highres: highres,
	}
}
