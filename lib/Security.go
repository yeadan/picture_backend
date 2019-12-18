package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"
	"fmt"
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/yeadan/picture_backend/api/data"
	"github.com/yeadan/picture_backend/api/models"
)


type JSONToken struct {
	Token string `json:"token"`
}

type TokenHMAC struct {
	UserID int
}

// ValidateCredent Buscar si existe un usuario para estos credenciales
func ValidateCredent(cred *models.Credentials, db *gorm.DB) *models.User {
	user := &models.User{}

	passHash := sha256.New()
	passHash.Write([]byte(cred.Password))
	newpass := fmt.Sprintf("%x", passHash.Sum(nil))

	db.Where("username = ? AND password = ?", cred.Username, newpass).Find(user)
	if user.UserID > 0 {
		return user
	} else {
		return nil
	}
}

// CreateHMAC - Crea un token para un usuario
func CreateHMAC(usr *models.User, cacheClient data.CacheProvider) (token string, err error) {
	h := hmac.New(sha256.New, []byte(os.Getenv("secret_key")))
	_, err = h.Write([]byte(usr.Username))
	if err == nil {
		token = hex.EncodeToString(h.Sum(nil))
		cacheClient.SetExpiration(token, usr, time.Hour)
	}
	return
}

// GetUserTokenCache - Devuelve el user del token
func GetUserTokenCache(tokenString string, cacheClient data.CacheProvider) *models.User {
	if validUser, exists := cacheClient.Get(tokenString); exists && validUser != nil {
		return validUser.(*models.User)
	}
	return nil
}

// UserAllowed - Permisos de roles y usuarios. Prevalece el id de usuario sobre el role
//usuario (es decir el allowUser) nil y role nil - sin restricciones
//usuario no nil y role nil - solo ese usuario
//usuario no nil y role no nil - ese usuario + ese role
//usuario nil y role no nil - ese role
//Si deja a un role, sea cual sea también dejará al role admin
func UserAllowed(user *models.User, allowUser *int, allowRole *string, w http.ResponseWriter) error {
	if user != nil { // Si existe un usuario
		if allowUser != nil { //Si se restringe a un ID d usuario
			// Se requiere un ID de usuario concreto o que haya un rol concreto
			if user.UserID == *allowUser || (allowRole != nil && *allowRole == user.Role) {
				return nil
			}
		} else {
			// No se requiere un usuario concreto, depende del rol
			if allowRole != nil { //Hay rol requerido 
				if user.Role == "admin" || *allowRole == user.Role {
					return nil
				}
			} else { //No hay restricciones
				return nil
			}
		}
	} 
	return errors.New("http error: Forbidden")
}

// GetString devuelve un puntero a un string
func GetString(text string) *string {
	return &text
}