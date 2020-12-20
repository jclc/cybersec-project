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
)

type registeredHandler struct {
	Pattern string
	Handler http.HandlerFunc
}

var handlers = make([]registeredHandler, 0)

func RegisterHandler(pattern string, handler http.HandlerFunc) {
	handlers = append(handlers, registeredHandler{pattern, handler})
}

// func StartFileServer(stop chan struct{})

func StartServer(port int) {
	if err := initTemplates(); err != nil {
		log.Println("Error parsing templates:", err)
		return
	}
	r := mux.NewRouter()
	for i := range handlers {
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
