package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/asaskevich/govalidator"
)

// Credentials para guardar datos de autenticación/autorización
type Credentials struct {
	Username string `gorm:"unique" json:"username" valid:"alphanum,length(3|15)"`
	Password string `json:"password" valid:"length(6|30)"`
}

// User - Struct con los datos del perfil del usuario
// Todas serán not null menos el avatar
type User struct {
	UserID int `gorm:"primary_key;not null" json:"user_id"`
	Credentials
	Role        string    `json:"role" valid:"alphanum,length(3|10)"`
	FullName    string    `json:"full_name" valid:"length(3|30)"`
	Registered  time.Time `json:"registered"`
	NumPictures int       `json:"num_pictures"`
	NumComments int       `json:"num_comment"`
	NumLikes    int       `json:"num_likes"`
	Avatar      *int 	  `gorm:"type:int REFERENCES images(image_id) ON DELETE SET NULL" json:"avatar"`
}

func (user User) Valid() bool {
	ok, err := govalidator.ValidateStruct(user)
	if err != nil {
		return false
	}
	return ok
}

func (user Credentials) Valid() bool {
	ok, err := govalidator.ValidateStruct(user)
	if err != nil {
		return false
	}
	return ok
}

func GetUsers(db *gorm.DB) []User {
	var users []User
	db.Find(&users)
	return users
}
func GetUser(id int, db *gorm.DB) *User {
	user := new(User)
	db.Find(user, id)
	if user.UserID == id {
		return user
	}
	return nil
}
// GetUserName busca el primer usuario con ese nombre en la BD. Si no existe, devuelve nil
func GetUserName(name string, db *gorm.DB) *User {
	user := new(User)
	db.Where("username = ?", name).First(&user)
	if user.Username == name {
		return user
	}
	return nil
}

func CreateUser(user *User, db *gorm.DB) {
	db.Create(user)
}
func EditUser(editUser *User, db *gorm.DB) {
	db.Save(editUser)
}

// NewCredentialsJSON Para pasar de JSO del body a objeto Credentials
func NewCredentialsJSON(jsonBytes []byte) *Credentials {
	cred := new(Credentials)
	err := json.Unmarshal(jsonBytes, cred)
	if err == nil {
		return cred
	}
	return nil
}

// NewUserJSON Para pasar de JSON del body a objeto User
func NewUserJSON(jsonBytes []byte) *User {
	cred := new(User)
	err := json.Unmarshal(jsonBytes, cred)
	if err == nil {
		return cred
	}
	return nil
}
