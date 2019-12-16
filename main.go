package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RGlimadigital/Tareas-Go/data"
	"github.com/RGlimadigital/Tareas-Go/lib"
	"github.com/RGlimadigital/Tareas-Go/middlewares"
	"github.com/RGlimadigital/Tareas-Go/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	image_manger "github.com/graux/image-manager"
	"github.com/urfave/negroni"
)

func login(w http.ResponseWriter, r *http.Request) {

	jsonBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cred := models.NewCredentialsJSON(jsonBytes)
	if cred == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, _ := data.ConnectDB()
	defer db.Close()
	validUser := lib.ValidateCredent(cred, db)

	if validUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// token, err := lib.CreateJWT(validUser)
	token, err := lib.CreateToken(validUser, data.GetCacheClient())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(lib.JSONToken{Token: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) //Token
}

//Codigo copiado de login
func register(w http.ResponseWriter, r *http.Request) {

	jsonBytes, err := ioutil.ReadAll(r.Body)
	println("test")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := models.NewUserJSON(jsonBytes)
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, _ := data.ConnectDB()
	defer db.Close()
	if err := db.Create(user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(fmt.Sprintf("error creating user: %s", err))
		return
	}

	responseBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes) //Token
}

func validateUser(w http.ResponseWriter, r *http.Request) *models.User {
	auto := r.Header.Get("Authorization")
	if len(auto) > 0 && strings.Contains(auto, "Bearer ") {
		tokenString := strings.Split(auto, " ")[1]
		db, _ := data.ConnectDB()
		defer db.Close()
		userValid := lib.GetUserJWT(tokenString, db)
		if userValid != nil {
			return userValid
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	return nil
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	db, _ := data.ConnectDB()
	defer db.Close()
	jsonTasks, err := json.Marshal(models.GetTasks(db))
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonTasks)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	if idStr, ok := mux.Vars(r)["id"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		id, _ := strconv.Atoi(idStr)
		task := models.GetTask(id, db)
		if task != nil {
			jsonTask, err := json.Marshal(task)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonTask)
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

func createTask(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil { // Validar userValid is *models.User
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			task := new(models.Task)
			err := json.Unmarshal(jsonBytes, task)
			if err == nil && task.Valid() {
				task.UserID = userValid.(*models.User).UserID
				db, _ := data.ConnectDB()
				defer db.Close()
				models.AddTask(task, db)
				w.Header().Set("Location", fmt.Sprintf("/tasks/%d", task.TaskID))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				jsonTask, _ := json.Marshal(task)
				w.Write(jsonTask)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func editTask(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			task := models.GetTask(id, db)
			if task != nil {
				if task.UserID == userValid.(*models.User).UserID {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						editTask := new(models.Task)
						err := json.Unmarshal(jsonBytes, editTask)
						if err == nil && task.Valid() {
							editTask.UserID = userValid.(*models.User).UserID
							editTask.TaskID = task.TaskID
							models.EditTask(editTask, db)
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

func getPictures(w http.ResponseWriter, r *http.Request) {
	db, _ := data.ConnectDB()
	defer db.Close()
	jsonPictures, err := json.Marshal(models.GetPictures(db))
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonPictures)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func createImage(w http.ResponseWriter, r *http.Request) {
	var imgType = "picture" //TODO Get from request
	imagesPath, err := filepath.Abs("./public/images-files")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error creating images %s\n", err)
		return
	}
	imgManager := image_manger.NewImageManager(imagesPath)
	var uuids []uuid.UUID
	imageBytes, err := ioutil.ReadAll(r.Body)
	if imgType == "avatar" {
		uuids, err = imgManager.ProcessImageAsSquare(imageBytes)
	} else if imgType == "picture" {
		uuids, err = imgManager.ProcessImageAs16by9(imageBytes)
	}
	thumb := uuids[0].String() + ".jpg"
	lowres := uuids[1].String() + ".jpg"
	highres := uuids[2].String() + ".jpg"
	userValid := r.Context().Value(middlewares.UserKey).(*models.User)
	img := models.NewImage(userValid.UserID, thumb, lowres, highres)
	db, _ := data.ConnectDB()
	defer db.Close()

	if err = db.Create(img).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error creating images %s\n", err)
		return
	}
	w.Header().Add("Location", fmt.Sprintf("/images/%d", img.ID))
	w.WriteHeader(http.StatusCreated)
}

//websocket

func createPicture(w http.ResponseWriter, r *http.Request) {
	//Recibe el Array de Bytes do body
	jsonBytes, err := ioutil.ReadAll(r.Body)

	//picture recibe el Array de bytes
	picture := models.NewPictureJSON(jsonBytes)
	if picture == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Conectando con la base de datos
	db, _ := data.ConnectDB()
	defer db.Close()

	if err := db.Create(picture).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(fmt.Sprintf("error creating picture: %s", err))
		return
	}

	//Se ha ido todo bien
	responseBytes, err := json.Marshal(picture)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Creando o Header
	w.Header().Add("Content-Type", "aplication/json")
	w.Header().Add("Location", fmt.Sprintf("/pictures/%d", picture.ID))
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes) //Picture
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			task := models.GetTask(id, db)
			if task != nil {
				if task.UserID == userValid.(*models.User).UserID {
					models.DeleteTask(task, db)
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

func main() {

	//Inicializando Router y Rutas
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users/login", login).Methods(http.MethodPost)
	router.HandleFunc("/users", register).Methods(http.MethodPost)
	router.HandleFunc("/tasks", getTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id:[0-9]+}", getTask).Methods(http.MethodGet)
	router.HandleFunc("/tasks", createTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id:[0-9]+}", editTask).Methods(http.MethodPut)
	router.HandleFunc("/images", createImage).Methods(http.MethodPost)
	router.HandleFunc("/pictures", getPictures).Methods(http.MethodGet)
	router.HandleFunc("/pictures", createPicture).Methods(http.MethodPost)
	router.HandleFunc("/websocket", models.ConnectWebSocket).Methods(http.MethodGet)

	router.Use(middlewares.AuthUser)

	// Creamos Negroni y registramos middlewares
	middle := negroni.Classic()
	middle.UseHandler(router)

	data.InitDB()

	http.ListenAndServe(":8080", middle)
}
