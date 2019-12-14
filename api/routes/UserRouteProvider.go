package routes

import (
	"crypto/sha256"
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

// GetRoutesUsers contiene las rutas del "usuario"
func GetRoutesUsers(r *mux.Router) {
	r.HandleFunc("/users/login", login).Methods(http.MethodPost)
	r.HandleFunc("/users", signup).Methods(http.MethodPost)
	r.HandleFunc("/users", getAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id:[0-9]+}", getUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods(http.MethodPut)
}
// updateUser - Para modificar un usuario. El username, el id y la fecha de creación no se pueden modificar.
//El role solo puede modificarlo un admin
func updateUser(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			user := models.GetUser(id,db)
			if user != nil {
				jsonBytes, err := ioutil.ReadAll(r.Body)
				if err == nil {
					editUser := new(models.User)
					err := json.Unmarshal(jsonBytes, editUser)
					//Comprobamos que no haya error en el Unmarshall y la validación de los datos con govalidator
					if err == nil && editUser.Valid(){
						errAuth := lib.UserAllowed(userValid.(*models.User), &id, lib.GetString("admin"), w)
						if errAuth == nil {
							//Codificando el nuevo password
							passHash := sha256.New()
							passHash.Write([]byte(editUser.Password))
							editUser.Password = fmt.Sprintf("%x", passHash.Sum(nil))
							user.Password = editUser.Password
							user.FullName = editUser.FullName
							errAuth = lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
							if errAuth == nil {
								//Solo un admin podra cambiar el role
								user.Role = editUser.Role
							}
							models.EditUser(user, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusForbidden)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

//Ver un usuario. Los usuarios registrados pueden verlo
func getUser(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { //  *models.User
		if idStr, ok := mux.Vars(r)["id"]; ok {
			db, _ := data.ConnectDB()
			defer db.Close()
			id, _ := strconv.Atoi(idStr)
			user := models.GetUser(id, db)
			if user != nil {
				jsonTask, err := json.Marshal(user)
				if err == nil {
					//Comprobar permisos y roles
					errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("user"), w)
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

// getAllUsers - Listar todos los usuarios. SOLO admins
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		db, _ := data.ConnectDB()
		defer db.Close()
		jsonTasks, err := json.Marshal(models.GetUsers(db))
		if err == nil {
			// Tiene permisos el usuario?
			errAuth := lib.UserAllowed(userValid.(*models.User), nil, lib.GetString("admin"), w)
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

// signup es la función para registrar usuarios. Encripta el password con SHA256
func signup(w http.ResponseWriter, r *http.Request) {

	//Cojemos el body y miramos si no da fallos al leerlo
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Miramos que el usuario sea un json válido
	user := models.NewUserJSON(jsonBytes)
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Comprobamos validación con govalidator
	//Algunos campos requieren longitudes máximas y mínimas
	if !user.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Encriptamos el password
	passHash := sha256.New()
	passHash.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", passHash.Sum(nil))

	//Conectamos a la base de datos
	db, _ := data.ConnectDB()
	defer db.Close()

	// Miramos si existe el usuario
	//La comprobación es doble, gorm también la haría por el unique
	//de username, pero así devolvemos el status y el mensaje 
	if models.GetUserName(user.Username, db) != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("El usuario ya existe"))
		return
	}
	//Usuario creado
	models.CreateUser(user, db)

	responseBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Lo devolvemos en JSON
	w.Header().Add("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}

//login sirve para logearse, creando un token HMAC para el usuario
func login(w http.ResponseWriter, r *http.Request) {
	//Cojemos el body y miramos si no da fallos al leerlo
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Miramos que la credencial sea un json válido
	cred := models.NewCredentialsJSON(jsonBytes)
	if cred == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Validamos el user y pass con govalidator
	//En principio no haría falta porque se hace en el signup
	if !cred.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Conectamos con la base de datos
	db, _ := data.ConnectDB()
	defer db.Close()
	//Comprobamos el user y pass en la base de datos
	validUser := lib.ValidateCredent(cred, db)
	if validUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//Creamos token HMAC para el usuario y lo ponemos en cache
	token, err := lib.CreateHMAC(validUser, data.GetCacheClient())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Cambiado por si se utilizaba el front de public
	//responseBytes, err := json.Marshal(lib.JSONToken{Token: token},)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	type Result struct {
		UserID int `json:"user_id"`
		Role string `json:"role"`
		Token string `json:"token"`
	}

	dat := Result{UserID: validUser.UserID,Role: validUser.Role, Token: token}
	responseBytes, err := json.Marshal(dat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	return
	}

	w.Header().Add("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) 
}
