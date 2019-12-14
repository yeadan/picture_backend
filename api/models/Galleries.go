package models

import (
	"github.com/jinzhu/gorm"
)

// Estructura de las galerías. Si borro un usuario, se borran todas sus galerías
type Gallery struct {
	GalleryID int `gorm:"primary_key" json:"gallery_id"`
	UserID    int `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE CASCADE" json:"user_id"`
	Title     string `json:"title" `
}
//Estructura de la relación entre Picture y Gallery
type PicGallery  struct{
	PicGalleryID int `gorm:"primary_key" json:"pic_gallery_id"`
	PictureID int `gorm:"not null;type:int REFERENCES pictures(picture_id) ON DELETE CASCADE" json:"picture_id"`
	GalleryID int `gorm:"not null;type:int REFERENCES galleries(gallery_id) ON DELETE CASCADE" json:"gallery_id"`
}

//coger galería por ID
func GetGallery(id int, db *gorm.DB) *Gallery {
	gallery := new(Gallery)
	db.Find(gallery, id)
	if gallery.GalleryID == id {
		return gallery
	}
	return nil
}
//coger picgallery por ID
func GetPicGallery(id int, db *gorm.DB) *PicGallery {
	picgallery := new(PicGallery)
	db.Find(picgallery, id)
	if picgallery.PicGalleryID == id {
		return picgallery
	}
	return nil
}
// GetPic2Gallery - Devuelve una picgallery entrando el id de una galeria y de una picture
//Si no existe, devuelve null
func GetPic2Gallery(id int, pic int, db *gorm.DB) *PicGallery {
	picgallery := new(PicGallery)
	db.Where("gallery_id = ? AND picture_id = ?", id,pic).Find(picgallery)
	if picgallery.GalleryID == id {
		return picgallery
	}
	return nil
}

// Mira si ya existe una picgallery
func ExistPicGallery(gal int, photo int, db *gorm.DB) bool {
	pic := new(PicGallery)
	result := db.Where("gallery_id = ? AND picture_id = ?", gal,photo).First(&pic)
	if result.RecordNotFound() {
		return false
	}
	return true
}


//GetPicsGallery - Lista todas las pictures de una galería
func GetPicsGallery(id int, db *gorm.DB) []*PicGallery {
	 pics := []*PicGallery{}
	db.Where("gallery_id = ?", id).Find(&pics)
	return pics
}

func CreateGallery(gallery *Gallery, db *gorm.DB) {
	db.Create(gallery)
}
func DeleteGallery(gallery *Gallery, db *gorm.DB) {
	db.Delete(gallery)
}

func EditGallery(gallery *Gallery, db *gorm.DB) {
	db.Save(gallery)
}
func CreatePicGallery(picgallery *PicGallery, db *gorm.DB) {
	db.Create(picgallery)
}
func DeletePicGallery(picgallery *PicGallery, db *gorm.DB) {
	db.Delete(picgallery)
}
//no hace falta
func EditPicGallery(picgallery *PicGallery, db *gorm.DB) {
	db.Save(picgallery)
}