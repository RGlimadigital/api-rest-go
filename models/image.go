package models

import "github.com/jinzhu/gorm"

type Image struct {
	gorm.Model
	UserID  int    `json:"user_id"`
	Thumb   string `json:"thumb_url"`
	Lowres  string `json:"lowres_url"`
	HighRes string `json:"highRes_url"`
}

func NewImage(userID int, thumb string, lowres string, highRes string) *Image {
	return &Image{
		UserID:  userID,
		Thumb:   thumb,
		Lowres:  lowres,
		HighRes: highRes,
	}
}
