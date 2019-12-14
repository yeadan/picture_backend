package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

//Like struct es la estructura de los like
type Like struct {
	LikeID 	int `gorm:"primary_key" json:"like_id"`
	UserID    int `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE CASCADE" json:"user_id"`
	PictureID int `gorm:"not null;type:int REFERENCES pictures(picture_id) ON DELETE CASCADE" json:"picture_id"`
	Liked     time.Time `json:"liked"`
}
func GetLike(id int, db *gorm.DB) *Like {
	like := new(Like)
	db.Find(like, id)
	if like.LikeID == id {
		return like
	}
	return nil
}
//Mira si para esa foto el usuario ya tiene un like, ya que no se acumulan
func ExistLike(user int, photo int, db *gorm.DB) bool {
	like := new(Like)
	result := db.Where("user_id = ? AND picture_id = ?", user,photo).First(&like)
	if result.RecordNotFound() {
		return false
	}
	return true
}
func CreateLike(like *Like, db *gorm.DB) {
	db.Create(like)
}
func DeleteLike(del *Like, db *gorm.DB) {
	db.Delete(del)
}
//no hace falta
func EditLike(editLike *Like, db *gorm.DB) {
	db.Save(editLike)
}