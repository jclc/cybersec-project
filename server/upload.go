package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/upload/{id:[1-9][0-9]*}/", handleUploadCRUD)
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
		upload.Delete()
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
