package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)
//Image struct es la estructura de las imágenes
type Image struct {
	ImageID    int    `gorm:"primary_key;not null" json:"image_id"`
	ThumbURL   string `gorm:"not null" json:"thumb_url"`
	LowResURL  string `gorm:"not null" json:"low_res_url"`
	HighResURL string `gorm:"not null" json:"high_res_url"`
}

//Función automática del gorm para mostrar los paths de las imágenes completos
func (u *Image) AfterFind() (err error) {
	u.ThumbURL = fmt.Sprintf("/images/%s.jpg", u.ThumbURL)
  	u.LowResURL = fmt.Sprintf("/images/%s.jpg", u.LowResURL)
	u.HighResURL = fmt.Sprintf("/images/%s.jpg", u.HighResURL)
   return
}

//Coger imagen por ID
func GetImage(id int, db *gorm.DB) *Image {
	image := new(Image)
	db.Find(image, id)
	if image.ImageID == id {
		return image
	}
	return nil
}
func CreateImage(image *Image, db *gorm.DB) {
	db.Create(image)
}
func DeleteImage(del *Image, db *gorm.DB) {
	db.Delete(del)
}
func EditImage(editImage *Image, db *gorm.DB) {
	db.Save(editImage)
}
//Coger la última imagen añadida, actualmente no utilizado
func LastImage(db *gorm.DB) *Image{
	 lastImage := new(Image)
	db.Last(lastImage)
	return lastImage
}