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

// GetRoutesGallerys contiene las rutas de las galerías/picgaleries
func GetRoutesGallery(r *mux.Router) {
	r.HandleFunc("/gallery", createGallery).Methods(http.MethodPost)
	r.HandleFunc("/gallery/{id:[0-9]+}", deleteGallery).Methods(http.MethodDelete)
	r.HandleFunc("/gallery/{gal:[0-9]+}/{pic:[0-9]+}", deletePicGallery).Methods(http.MethodDelete)
	r.HandleFunc("/gallery/{gal:[0-9]+}/{pic:[0-9]+}", createPicGallery).Methods(http.MethodPost)
	r.HandleFunc("/gallery/{id:[0-9]+}", editGallery).Methods(http.MethodPut)
	r.HandleFunc("/gallery/{id:[0-9]+}", getGallery).Methods(http.MethodGet)
}

//getGallery - Devuelve todas las imágenes de una galería 
func getGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			gallery := models.GetGallery(id, db)
			if gallery != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &gallery.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					galleryJSON := []*models.Picture{}
					getPic := models.GetPicsGallery(id, db)
					for _,proba := range getPic {
						single := models.GetPicture(proba.PictureID,db)
						galleryJSON = append(galleryJSON,single)
					}
					jsonFinal, _ := json.Marshal(galleryJSON)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(jsonFinal)
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


// editGallery - Para cambiar el título de la galería
func editGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			gallery := models.GetGallery(id, db)
			if gallery != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &gallery.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						edited := new(models.Gallery)
						err := json.Unmarshal(jsonBytes, edited)
						if err == nil {
							gallery.Title = edited.Title
							models.EditGallery(gallery, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusBadRequest)
							w.Write([]byte("El body no contiene un JSON válido"))
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("Error leyendo el body"))
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

// createPicGallery - Mete una picture en una galería. No puede poner fotos repetidas dentro de la galería
func createPicGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if galStr, ok := mux.Vars(r)["gal"]; ok {
			if picStr, ok := mux.Vars(r)["pic"]; ok {
				gal, _ := strconv.Atoi(galStr)
				pic, _ := strconv.Atoi(picStr)
				db, _ := data.ConnectDB()
				defer db.Close()
				gallery := models.GetGallery(gal,db)
				picture := models.GetPicture(pic,db)
				if gallery != nil && picture != nil {
					errAuth := lib.UserAllowed(userValid.(*models.User), &gallery.UserID, nil, w)
					if errAuth == nil && gallery.UserID == picture.UserID{ //Solo el propietario de foto y galería puede añadir
						if models.ExistPicGallery(gal,pic,db) {
							w.WriteHeader(http.StatusBadRequest)
							w.Write([]byte("No se puede volver a añadir la misma foto a la galería"))
							return
						}
						picgallery := new(models.PicGallery)
						picgallery.PictureID = picture.PictureID
						picgallery.GalleryID = gallery.GalleryID
						models.CreatePicGallery(picgallery, db)
						w.Header().Set("Location", fmt.Sprintf("/gallery/%d/%d", picgallery.GalleryID,picgallery.PictureID))
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						jsonTask, _ := json.Marshal(picgallery)
						w.Write(jsonTask)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}


// deletePicGallery - Se envía el id de la galería y de la foto en esa galería que se quiere borrar
//y se borra la PicGallery resultante
func deletePicGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if galStr, ok := mux.Vars(r)["gal"]; ok {
			if picStr, ok := mux.Vars(r)["pic"]; ok {
				gal, _ := strconv.Atoi(galStr)
				pic, _ := strconv.Atoi(picStr)
				db, _ := data.ConnectDB()
				defer db.Close()
				pic_gallery := models.GetPic2Gallery(gal,pic, db)
				if pic_gallery != nil {
					gallery := models.GetGallery(gal,db)
					errAuth := lib.UserAllowed(userValid.(*models.User), &gallery.UserID, lib.GetString("admin"), w)
					if errAuth == nil {
						models.DeletePicGallery(pic_gallery, db)
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
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}


// deleteGallery - Elimina una gallery 
func deleteGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			gallery := models.GetGallery(id, db)
			if gallery != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &gallery.UserID, nil, w)
				if errAuth == nil {
					// Se borrarán todas las fotos de la galería (las picgalleries) por el on delete cascade
					models.DeleteGallery(gallery, db)
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
// createGallery - Creamos el gallery 
func createGallery(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { 
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			gallery := new(models.Gallery)
			err := json.Unmarshal(jsonBytes, gallery)
			if err == nil { 
				errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
				if errAuth == nil {
					gallery.UserID = userValid.(*models.User).UserID
					db, _ := data.ConnectDB()
					defer db.Close()			
					models.CreateGallery(gallery, db)
					w.Header().Set("Location", fmt.Sprintf("/gallery/%d", gallery.GalleryID))
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					jsonTask, _ := json.Marshal(gallery)
					w.Write(jsonTask)
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

