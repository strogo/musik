package main

import (
	"io"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/imthaghost/musik/soundcloud"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

var songlist = [...]string{"https://soundcloud.com/polo-g/polo-g-feat-juice-wrld-flex", "https://soundcloud.com/roddyricch/the-box", "https://soundcloud.com/lil-baby-4pf/sum-2-prove"}
var old string

func main() {
	// random integer

	e := echo.New()

	// Log Output
	e.Use(middleware.Logger())
	// stream
	e.Use(middleware.Recover())
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// static files
	e.Static("/", "assets")
	// template render
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	// Named route "index"
	e.GET("/", func(c echo.Context) error {
		rand.Seed(time.Now().Unix())
		var n = rand.Int() % len(songlist)

		// err := os.Remove("assets/music/" + old)
		// if err != nil {

		// }
		songname, image, path := soundcloud.ExtractSong(songlist[n])
		old = path
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{

			"name":    songname,
			"artwork": image,
			"song":    "music/" + path,
		})
	}).Name = "index"

	e.Logger.Fatal(e.Start(":5000"))
}
