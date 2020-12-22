package server

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template" // there's html/template, but I don't know the difference so I'll use this
	"time"

	"github.com/jclc/cybersec-project/database"
)

const templatePath = `./templates/`

// BaseContext contains data for rendering the base template, such as current
// page's location on the navigation bar or the currently logged in user (if any).
type BaseContext struct {
	Title string        // Window title, can be empty
	User  database.User // Currently logged in user. Zeroed if logged out.
	Nav   string        // Current sub-page indicated by the navigation bar
}

// NavItems is a list of the links available in the navigation menu.
// The first string in each tuple is the identifier.
// The second string is the lable to be shown on the link.
// The third string is the URL that the link should refer to.
var navItems = [][3]string{
	{"index", "Home", "/"},
}

// Appended to navItems when user is logged in.
var navItemsLoggedIn = [][3]string{
	{"files", "My files", "/myfiles/"},
	{"logout", "Log out", "/login/logout/"},
}

// Appended to navItems when user is not logged in.
var navItemsLoggedOut = [][3]string{
	{"login", "Log in/Register", "/login/"},
}

var tmpl map[string]*template.Template

var renderFuncs = template.FuncMap{
	"fdate": func(d time.Time) string {
		return d.Format("02.01.2006 15:04 -0700")
	},
}

func initTemplates() error {
	files, err := ioutil.ReadDir(templatePath)
	if err != nil {
		return err
	}
	tmpl = make(map[string]*template.Template)
	for _, f := range files {
		name := f.Name()
		if name == "base.html" {
			continue
		}
		t := template.New(name)
		// var t *template.Template
		t.Funcs(renderFuncs)
		t, err := t.ParseFiles(
			filepath.Join(templatePath, name),
			filepath.Join(templatePath, "base.html"),
		)
		if err != nil {
			return err
		}
		tmpl[name] = t
		log.Println("Template parsed:", name)
	}
	log.Println("Templates parsed")
	return nil
}

// RenderTemplate renders the `content` template inside the base template and
// writes it into w.
func RenderTemplate(w io.Writer, content string, base *BaseContext, data interface{}) {
	var ctx struct {
		NavItems    [][3]string
		BaseCtx     *BaseContext
		ContentData interface{}
	}
	ctx.NavItems = navItems
	ctx.ContentData = data
	if base != nil {
		ctx.BaseCtx = base
	} else {
		ctx.BaseCtx = &BaseContext{}
	}
	if ctx.BaseCtx.User.ID == 0 {
		ctx.NavItems = append(ctx.NavItems, navItemsLoggedOut...)
	} else {
		ctx.NavItems = append(ctx.NavItems, navItemsLoggedIn...)
	}

	err := tmpl[content].ExecuteTemplate(w, "base", &ctx)
	if err != nil {
		log.Printf("Error in template %s: %v\n", content, err)
	}
}
