package server

import (
	"net/http"
	"time"
)

func init() {
	RegisterHandler("/", handleIndex)
	RegisterHandler("/users", handleUsers)
	// RegisterHandler("/storage", handleTemp)
}

type Upload struct {
	Date       time.Time
	Filename   string
	Visibility string
	Action     string
}

// HandleIndex displays the front page to the user.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := struct {
		Uploads []Upload
	}{
		Uploads: []Upload{
			{
				Date:       time.Now(),
				Filename:   "a.jpg",
				Visibility: "Public",
				Action:     "Delete",
			},
			{
				Date:       time.Now().Add(time.Hour * -531),
				Filename:   "b.jpg",
				Visibility: "Public",
				Action:     "Delete",
			},
			{
				Date:       time.Now().Add(time.Hour * -5312),
				Filename:   "c.jpg",
				Visibility: "Hidden",
				Action:     "Delete",
			},
		},
	}
	RenderTemplate(w, "index.html", nil, &ctx)
}

// func handleTemp(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		q := r.URL.Query()
// 		file, ok := q["file"]
// 		err := storage.Get(w, r, "image.jpg", "image")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func handleUsers(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "list_users.html", nil, []string{
		"Admin",
		"Superuser",
		"Literal God",
		"Noob",
	})
}
