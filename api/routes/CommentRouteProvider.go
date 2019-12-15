package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yeadan/proyect-image/api/data"
	"github.com/yeadan/proyect-image/api/middlewares"
	"github.com/yeadan/proyect-image/api/models"
	"github.com/yeadan/proyect-image/lib"
)

// GetRoutesComments contiene las rutas del "comment"
func GetRoutesComments(r *mux.Router) {
	r.HandleFunc("/comment", createComment).Methods(http.MethodPost)
	r.HandleFunc("/comment/{id:[0-9]+}", editComment).Methods(http.MethodPut)
	r.HandleFunc("/comment/{id:[0-9]+}", deleteComment).Methods(http.MethodDelete)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			comment := models.GetComment(id, db)
			if comment != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &comment.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					models.DeleteComment(comment, db)
					picture := models.GetPictureID(comment.PictureID,db)
					picture.NumComments--
					models.EditPicture(picture,db)
					user := models.GetUser(comment.UserID,db)
					user.NumComments--
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
func createComment(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			comment := new(models.Comment)
			err := json.Unmarshal(jsonBytes, comment)
			if err == nil { //&& task.Valid() {
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					comment.UserID = userValid.(*models.User).UserID
					db, _ := data.ConnectDB()
					defer db.Close()
					//aumentamos en 1 el número de comentarios de picture y user
					picture := models.GetPictureID(comment.PictureID,db)
					if picture != nil {
						models.CreateComment(comment, db)
						picture.NumComments++
						models.EditPicture(picture,db)
						user := models.GetUser(comment.UserID,db)
						user.NumComments++
						models.EditUser(user,db)
						w.Header().Set("Location", fmt.Sprintf("/comment/%d", comment.CommentID))
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						jsonTask, _ := json.Marshal(comment)
						w.Write(jsonTask)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func editComment(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			comment := models.GetComment(id, db)
			if comment != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &comment.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						edited := new(models.Comment)
						err := json.Unmarshal(jsonBytes, edited)
						if err == nil { //&& task.Valid() {
							comment.Comment = edited.Comment
							models.EditComment(comment, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
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
