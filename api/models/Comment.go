package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

//Comment struct de los comentarios
type Comment struct {
	CommentID int       `gorm:"primary_key" json:"comment_id"`
	UserID    int       `gorm:"not null;type:int REFERENCES users(user_id) ON DELETE CASCADE" json:"user_id"`
	PictureID int       `gorm:"not null;type:int REFERENCES pictures(picture_id) ON DELETE CASCADE" json:"picture_id"`
	Created   time.Time `json:"created"`
	Comment   string    `json:"comment"`
}

// NewComment constructor para crear nuevos comentarios
func NewComment(commentID int, userID int, pictureID int, comment string) *Comment {
	return &Comment{
		commentID,
		userID,
		pictureID,
		time.Now(),
		comment,
	}
}
func GetComment(id int, db *gorm.DB) *Comment {
	comment := new(Comment)
	db.Find(comment, id)
	if comment.CommentID == id {
		return comment
	}
	return nil
}
func CreateComment(comment *Comment, db *gorm.DB) {
	db.Create(comment)
}
func DeleteComment(del *Comment, db *gorm.DB) {
	db.Delete(del)
}
func EditComment(editComment *Comment, db *gorm.DB) {
	db.Save(editComment)
}