package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type registeredHandler struct {
	Pattern string
	Handler http.HandlerFunc
}

var handlers = make([]registeredHandler, 0)

func RegisterHandler(pattern string, handler http.HandlerFunc) {
	handlers = append(handlers, registeredHandler{pattern, handler})
}

var store *sessions.CookieStore

func StartServer(port int, sessionKey string) {
	if err := initTemplates(); err != nil {
		log.Println("Error parsing templates:", err)
		return
	}

	store = sessions.NewCookieStore([]byte(sessionKey))

	r := mux.NewRouter()
	r.StrictSlash(true)
	for i := range handlers {
		log.Println("Registering handler for", handlers[i])
		r.HandleFunc(handlers[i].Pattern, handlers[i].Handler)
	}
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: r,
	}
	log.Println("Initalising server")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt) // Catch ^C
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("HTTP server successfully shut down")
}
