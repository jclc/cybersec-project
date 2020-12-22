package server

import (
	"fmt"
	"net/http"

	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/", handleIndex)
	RegisterHandler("/myfiles/", handleFiles)
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	// log.Println("/myfiles/")
	user := getRegisteredUser(r)
	if user.ID == 0 {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/user/%d/", user.ID), http.StatusFound)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// log.Println("/")
	ctx := struct {
		Users []database.User
	}{}
	ctx.Users = database.GetUsers()

	baseCtx := BaseContext{
		Nav:  "index",
		User: getRegisteredUser(r),
	}
	RenderTemplate(w, "index.html", &baseCtx, &ctx)
}
