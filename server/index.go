package server

import "net/http"

func init() {
	RegisterHandler("/", HandleIndex)
}

// HandleIndex displays the front page to the user.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
