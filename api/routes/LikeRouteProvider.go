package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yeadan/picture_backend/api/data"
	"github.com/yeadan/picture_backend/api/middlewares"
	"github.com/yeadan/picture_backend/api/models"
	"github.com/yeadan/picture_backend/lib"
)

// GetRoutesLikes contiene las rutas del "like"
func GetRoutesLikes(r *mux.Router) {
	r.HandleFunc("/like", createLike).Methods(http.MethodPost)
	r.HandleFunc("/like/{id:[0-9]+}", deleteLike).Methods(http.MethodDelete)
}

// deleteLike - Es la función de unlike. Elimina un like si existe, y actualiza los valores de numLikes en user y picture
func deleteLike(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			like := models.GetLike(id, db)
			if like != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &like.UserID, nil, w)
				if errAuth == nil {
					models.DeleteLike(like, db)
					picture := models.GetPictureID(like.PictureID,db)
					picture.NumLikes--
					models.EditPicture(picture,db)
					user := models.GetUser(like.UserID,db)
					user.NumLikes--
					models.EditUser(user,db)
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
// createLike - Creamos el like y actualizamos el user y la picture
func createLike(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			like := new(models.Like)
			err := json.Unmarshal(jsonBytes, like)
			if err == nil { 
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					like.UserID = userValid.(*models.User).UserID
					db, _ := data.ConnectDB()
					defer db.Close()			
					picture := models.GetPictureID(like.PictureID,db)
					if picture != nil {
						if models.ExistLike(like.UserID,like.PictureID,db) {
							w.WriteHeader(http.StatusBadRequest)
							w.Write([]byte("Ya existe el like de ese usuario en esa foto"))
							return
						}
						models.CreateLike(like, db)
						picture.NumLikes++
						models.EditPicture(picture,db)
						user := models.GetUser(like.UserID,db)
						user.NumLikes++
						models.EditUser(user,db)	
						w.Header().Set("Location", fmt.Sprintf("/like/%d", like.LikeID))
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						jsonTask, _ := json.Marshal(like)
						w.Write(jsonTask)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("El body no contiene un JSON válido"))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error leyendo el body"))
		}
	}
}

