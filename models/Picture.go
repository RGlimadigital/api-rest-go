package models

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Picture struct {
	gorm.Model
	ImageID      int    `json:"image_id"`
	Title        string `json: "title"`
	Descritption string `json: "description"`
	Image        Image
}

func newPicture(imageID int, title string, description string) *Picture {
	return &Picture{
		ImageID:      imageID,
		Title:        title,
		Descritption: description,
	}
}

//Comentaando NewPictureJSON me devolvendo un *Picture
func NewPictureJSON(jsonBytes []byte) *Picture {
	picture := new(Picture)
	err := json.Unmarshal(jsonBytes, picture)
	if err == nil {
		return picture
	}
	return nil
}

func GetPictures(db *gorm.DB) []Picture {
	var pictures []Picture
	db.Preload("Image").Preload("User").Order("created_at desc").Find(&pictures) //latest to be on top
	fmt.Println(pictures)
	return pictures
}
