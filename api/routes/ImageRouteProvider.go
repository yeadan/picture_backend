package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	image_manger "github.com/graux/image-manager"
	"github.com/yeadan/proyect-image/api/data"
	"github.com/yeadan/proyect-image/api/middlewares"
	"github.com/yeadan/proyect-image/api/models"
	"github.com/yeadan/proyect-image/lib"

	"github.com/google/uuid"
)

// GetRoutesImages contiene las rutas de picture e image
func GetRoutesImages(r *mux.Router) {
	r.HandleFunc("/images", createImage).Methods(http.MethodPost)
	r.HandleFunc("/images", getAllPictures).Methods(http.MethodGet)
	r.HandleFunc("/images/{id:[0-9]+}", editImage).Methods(http.MethodPut)
	r.HandleFunc("/images/{id:[0-9]+}", getPicturesUser).Methods(http.MethodGet)
	r.HandleFunc("/image/{id:[0-9]+}", getPicture).Methods(http.MethodGet)
	r.HandleFunc("/images/{id:[0-9]+}", deleteImage).Methods(http.MethodDelete)
	r.HandleFunc("/images/avatar", createImageAvatar).Methods(http.MethodPost)
}

//getAllPictures - Devuelve todas las imágenes ordenadas por recientes. Sin detalles del usuario
func getAllPictures(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		db, _ := data.ConnectDB()
		defer db.Close()
		jsonTasks, err := json.Marshal(models.GetPictures(db))
		if err == nil {
			// Tiene permisos el usuario?
			errAuth := lib.UserAllowed(userValid.(*models.User), nil, nil, w)
			if errAuth == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonTasks)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// deleteImage - Borra la imagen, el picture y decrementa en user el número de pictures
//Solo podrá borrarla el propietario o un admin
func deleteImage(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			picture := models.GetPicture(id, db)
			if picture != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &picture.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					// Ahora borramos la picture y la image asociada, pero antes miramos cuantos likes y comments tiene
					//para borrarlos en user. Los comments y likes se borran automáticamente al borrar la picture por
					//la defnición -> on delete cascade
					totComments := picture.NumComments
					totLikes := picture.NumLikes
					image := models.GetImage(picture.ImageID, db)
					models.DeletePicture(picture, db)
					if image != nil { // Si ya era nil, que no debería, ya no existe así que no devolvemos un error
						err := os.Remove(fmt.Sprintf("images/%s.jpg", image.ThumbURL))
						if err != nil {
							os.Remove(fmt.Sprintf("images/avatars/%s.jpg", image.ThumbURL))
							os.Remove(fmt.Sprintf("images/avatars/%s.jpg", image.LowResURL))
							os.Remove(fmt.Sprintf("images/avatars/%s.jpg", image.HighResURL))
							// Al borrar el avatar de un usuario, pasará a valer null el campo
							//avatar del user, por el on delete set null de su definición
						} else {
							os.Remove(fmt.Sprintf("images/%s.jpg", image.LowResURL))
							os.Remove(fmt.Sprintf("images/%s.jpg", image.HighResURL))
						}
						models.DeleteImage(image, db)
					}
					//Actualizamos el numero de pictures del usuario, el de likes y el de comments
					user := models.GetUser(picture.UserID, db)
					user.NumPictures--
					user.NumLikes = user.NumLikes - totLikes
					user.NumComments = user.NumComments - totComments
					models.EditUser(user, db)
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

// createImageAvatar - Crea un avatar para un usuario. Le pone de título 'avatar', aunque podrá editarlo
//Se crearán en el subdirectorio avatars dentro de images.
func createImageAvatar(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
			if errAuth == nil {
				db, _ := data.ConnectDB()
				imagesPath := "./images/avatars"
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				imgManager := image_manger.NewImageManager(imagesPath)
				var uuids []uuid.UUID
				uuids, err = imgManager.ProcessImageAsSquare(jsonBytes)
				if err == nil {
					//Creamos image
					image := new(models.Image)
					image.ThumbURL = fmt.Sprintf("avatars/%s", uuids[0])
					image.LowResURL = fmt.Sprintf("avatars/%s", uuids[1])
					image.HighResURL = fmt.Sprintf("avatars/%s", uuids[2])
					defer db.Close()
					models.CreateImage(image, db)
					//Creamos picture
					picture := new(models.Picture)
					picture.UserID = userValid.(*models.User).UserID
					picture.Image = *image
					picture.ImageID = image.ImageID
					picture.Created = time.Now()
					picture.Title = "avatar"
					//actualizamos el user
					user := models.GetUser(picture.UserID, db)
					user.NumPictures++
					user.Avatar = &image.ImageID
					models.EditUser(user, db)
					picture.User = models.GetUser(userValid.(*models.User).UserID, db)
					models.CreatePicture(picture, db)
					//Devolveremos por json la picture creada y el usuario actualizado, para ver que se ha creado bien
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					jsonUse, _ := json.Marshal(picture)
					w.Write(jsonUse)
					jsonUse, _ = json.Marshal(user)
					w.Write(jsonUse)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// getPicturesUser - Busca las pictures de un usuario, sin detalles del usuario
func getPicturesUser(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { //  *models.User
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			user := models.GetUser(id, db)
			if user != nil {
				jsonUser, err := json.Marshal(models.GetPicturesUser(id, db))
				if err == nil {
					errAuth := lib.UserAllowed(userValid.(*models.User), nil, nil, w)
					if errAuth == nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(jsonUser)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// createImage - Crea una image y una picture nuevas. Se le pasa una foto en formato jpeg en el body
//También actualiza el número de pictures del usuario
func createImage(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
			if errAuth == nil {
				db, _ := data.ConnectDB()
				imagesPath := "./images"
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				imgManager := image_manger.NewImageManager(imagesPath)
				var uuids []uuid.UUID
				uuids, err = imgManager.ProcessImageAs16by9(jsonBytes)
				if err == nil {
					image := new(models.Image)
					image.ThumbURL = fmt.Sprint(uuids[0])
					image.LowResURL = fmt.Sprint(uuids[1])
					image.HighResURL = fmt.Sprint(uuids[2])
					defer db.Close()
					models.CreateImage(image, db)
					//imagen añadida, que es la que hemos creado con CreateImage
					picture := new(models.Picture)
					picture.UserID = userValid.(*models.User).UserID
					picture.Image = *image
					picture.ImageID = image.ImageID
					picture.Created = time.Now()
					user := models.GetUser(picture.UserID, db)
					user.NumPictures++
					models.EditUser(user, db)
					picture.User = models.GetUser(userValid.(*models.User).UserID, db)
					models.CreatePicture(picture, db)
					w.Header().Set("Location", fmt.Sprintf("/image/%d", image.ImageID))
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					jsonUse, _ := json.Marshal(picture)
					w.Write(jsonUse)

				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

//No tiene sentido editar las images, así que solo se utilizará para editar títulos y descripciones del picture
func editImage(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			picture := models.GetPictureID(id, db)
			if picture != nil {
				errAuth := lib.UserAllowed(userValid.(*models.User), &picture.UserID, lib.GetString("admin"), w)
				if errAuth == nil {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						edited := new(models.Picture)
						err := json.Unmarshal(jsonBytes, edited)
						if err == nil {
							picture.Title = edited.Title
							picture.Description = edited.Description
							models.EditPicture(picture, db)
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

// getPicture - Enseña la picture con detalle (de usuario y de image)
func getPicture(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { //  *models.User
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			pic := models.GetPicture(id, db)
			if pic != nil {
				pic.User = models.GetUser(pic.UserID, db)
				jsonTask, err := json.Marshal(pic)
				if err == nil {
					errAuth := lib.UserAllowed(userValid.(*models.User), nil, nil, w)
					if errAuth == nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(jsonTask)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
