package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Picture Struct con los datos de una imagen.
type Picture struct {
	PictureID   int `gorm:"primary_key;not null" json:"picture_id"`
	UserID      int `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE CASCADE" json:"user_id"`
	ImageID     int `gorm:"not null;type:int REFERENCES images(image_id) ON DELETE CASCADE" json:"image_id"`
	Title       string `json:"title"`
	Description *string `json:"description"`
	Created     time.Time `json:"created"`
	NumLikes    int `json:"num_likes"`
	NumComments int	`json:"num_comments"`
	Image       Image `gorm:"foreignkey:ImageID" json:"image"` 
	User        *User `json:"user"` //no es FK por los errores del gorm con autoincrementos
}

//GetPictures - Lista todas las pictures, las más recientes primero
func GetPictures(db *gorm.DB) []*Picture {
	 pictures := []*Picture{}
	db.Order("created desc").Preload("Image").Find(&pictures)
	return pictures
}
// GetPicturesUser - Lista todas las pictures de un usuario concreto
func GetPicturesUser(user int, db *gorm.DB) []*Picture {
	pictures := []*Picture{}
	db.Where("pictures.user_id = ?", user).Preload("Image").Find(&pictures)
	return pictures
}

// GetPicture - Devuelve la picture con el ID
func GetPicture(id int, db *gorm.DB) *Picture {
	picture := new(Picture)
	db.Preload("Image").Find(picture, id)
	if picture.PictureID == id {
		return picture
	}
	return nil
}
func CreatePicture(picture *Picture, db *gorm.DB) {
	db.Create(picture)
}
func DeletePicture(del *Picture, db *gorm.DB) {
	db.Delete(del)
}
func EditPicture(editPicture *Picture, db *gorm.DB) {
	db.Save(editPicture)
}
