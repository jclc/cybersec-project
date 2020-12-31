package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jclc/cybersec-project/database"
	"github.com/jclc/cybersec-project/storage"
)

const (
	MaxUploadSize = 32 * 1024 * 1024
)

func init() {
	RegisterHandler("/upload/", handleUpload)
	RegisterHandler("/upload/{id:[1-9][0-9]*}/", handleUploadCRUD)
}

func storageName(id int64) string {
	return fmt.Sprintf("%d.bin", id)
}

// This function handles file uploads. It's super safe.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := getRegisteredUser(r)
	if user.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Error parsing upload:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	upload := database.Upload{
		Owner:      user.ID,
		Filename:   header.Filename,
		Timestamp:  time.Now().UTC(),
		Visibility: database.VHidden,
	}
	err = upload.Save()
	if err != nil {
		log.Println("Error saving upload to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	storageFile := storageName(upload.ID)
	err = storage.Store(storageFile, file)
	if err != nil {
		log.Println("Error storing file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/%d/", user.ID), http.StatusFound)
}

// This function handles all file operations except upload. It's super safe.
func handleUploadCRUD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	upload, err := database.UploadByID(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		err := storage.Get(w, r, storageName(upload.ID), upload.Filename)
		if err != nil {
			log.Println("Error getting upload:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PATCH":
		// Change visibility
		if upload.Visibility == database.VHidden {
			upload.Visibility = database.VPublic
		} else {
			upload.Visibility = database.VHidden
		}
		err := upload.Save()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "DELETE":
		storage.Delete(storageName(upload.ID))
		upload.Delete()
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
