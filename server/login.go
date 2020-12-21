package server

import (
	"log"
	"net/http"

	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/login", handleLoginPage)
	RegisterHandler("/login/auth", handleLoginAttempt)
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.html", nil, nil)
}

func handleLoginAttempt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var username, password string
	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "username":
			if len(v) > 0 {
				username = v[0]
			}
		case "password":
			if len(v) > 0 {
				password = v[0]
			}
		}
	}
	if password == "" || username == "" {
		RenderTemplate(w, "login.html", nil, "insufficient login information")
		return
	}

	// Keep track of logins just in case someone tries to do something nefarious.
	log.Printf("User '%s' with password '%s' attempted login\n", username, password)

	_, err := database.GetUser(username, password)
	if err != nil {
		RenderTemplate(w, "login.html", nil, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
