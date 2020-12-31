package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/user/{id:[1-9][0-9]*}/", handleUserPage)
	RegisterHandler("/user/{id:[1-9][0-9]*}/message/", handleMessage)
	RegisterHandler("/message/{id:[1-9][0-9]*}/", handleMessageCRUD)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	// log.Println("/user/id/")
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
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

	user, err := database.UserByID(id)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	data.User = user
	isOwnPage := baseCtx.User.ID == id
	data.LoggedInUser = baseCtx.User
	data.Uploads = database.GetUploads(id, !isOwnPage)
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

	var msg database.Message
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	msg.Recipient = id
	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "content":
			msg.Content = v[0]
		}
	}
	msg.Author = user.Username
	if msg.Content != "" {
		if err := msg.Save(); err != nil {
			log.Println("Error saving message:", err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/user/%d/", id), http.StatusFound)
}

func handleMessageCRUD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	msg, err := database.MessageByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte(err.Error()))
		return
	}

	// user := getRegisteredUser(r)
	// if user.Username != msg.Author && user.ID != msg.Recipient {
	// 	// Only let users delete their own messages or messages posted on their page.
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// TODO: We could add more operations than just DELETE
	switch r.Method {
	case "DELETE":
		msg.Delete()
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
