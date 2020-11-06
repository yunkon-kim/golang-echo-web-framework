package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

// Render renders a template document
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "aaa.html", nil)
}


func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello.html", map[string]interface{}{
		"name": "Dolly!",
	})
}


func main() {

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = t

	e.Static("/", "assets")
	e.Static("/js", "js")
	e.Static("/css", "css")
	e.File("/", "public/index.html")
	e.GET("/aaa", Index)
	e.GET("/hello", Hello)

	e.Logger.Fatal(e.Start(":8000"))
}
