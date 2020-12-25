package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/user/{id:[1-9][0-9]*}/", handleUserPage)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	// log.Println("/user/id/")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	baseCtx := BaseContext{
		Nav:  "index",
		User: getRegisteredUser(r),
	}

	data := struct {
		User         database.User
		Uploads      []database.Upload
		LoggedInUser database.User
		Messages     []database.Message
	}{}

	user, err := database.UserByID(int64(id))
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	data.User = user
	isOwnPage := baseCtx.User.ID == int64(id)
	data.LoggedInUser = baseCtx.User
	data.Uploads = database.GetUploads(int64(id), !isOwnPage)
	data.Messages = user.GetMessages()

	RenderTemplate(w, "user.html", &baseCtx, &data)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := getRegisteredUser(r)
	if user.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	var msg database.Message
	for k, v := range r.Form {
		switch k {
		case "content":
			msg.Content = v[0]
		case "recipient":
			msg.Content = v[0]
		}
	}
	msg.Author = user.Username
	msg.R
}
