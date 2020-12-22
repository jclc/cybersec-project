package server

import (
	"log"
	"net/http"

	"github.com/jclc/cybersec-project/database"
)

func init() {
	RegisterHandler("/login/logout/", handleLogout)
	RegisterHandler("/login/", handleLoginPage)
	RegisterHandler("/login/auth/", handleLoginAttempt)
	RegisterHandler("/login/register/", handleRegistration)
}

func getRegisteredUser(r *http.Request) database.User {
	session, _ := store.Get(r, "session-name")
	id, ok := session.Values["user_id"].(int64)
	if !ok || id == 0 {
		// not logged in
		log.Println("User is not logged in.")
		return database.User{}
	}
	username, _ := session.Values["username"].(string)
	log.Println("User", username, "is logged in.")
	return database.User{
		ID:       id,
		Username: username,
	}
}

func setRegisteredUser(w http.ResponseWriter, r *http.Request, user database.User) {
	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)
	log.Println("User's session has been set.")
}

func unsetRegisteredUser(w http.ResponseWriter, r *http.Request) {
	log.Println("User's session has been unset.")
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "user_id")
	delete(session.Values, "username")
	session.Save(r, w)
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	log.Println("/login/")
	baseCtx := BaseContext{
		Nav:  "login",
		User: getRegisteredUser(r),
	}
	RenderTemplate(w, "login.html", &baseCtx, nil)
}

func handleLoginAttempt(w http.ResponseWriter, r *http.Request) {
	log.Println("/login/auth/")
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

	user, err := database.GetUser(username, password)
	if err != nil {
		RenderTemplate(w, "login.html", nil, err.Error())
		return
	}

	setRegisteredUser(w, r, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	log.Println("/login/register/")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user database.User
	r.ParseForm()
	for k, v := range r.Form {
		switch k {
		case "username":
			if len(v) > 0 {
				user.Username = v[0]
			}
		case "password":
			if len(v) > 0 {
				user.Password = v[0]
			}
		case "social":
			if len(v) > 0 {
				user.SocialSecurity = v[0]
			}
		}
	}
	if user.Username == "" || user.Password == "" || user.SocialSecurity == "" {
		RenderTemplate(w, "login.html", nil, "insufficient user information")
		return
	}

	err := user.Create()
	if err != nil {
		RenderTemplate(w, "login.html", nil, err.Error())
		return
	}

	setRegisteredUser(w, r, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("/login/logout/")
	// log.Println("HELLO????")
	unsetRegisteredUser(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
