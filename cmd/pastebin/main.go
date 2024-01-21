package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	e.GET("/", index)
	e.GET("/paste/:id", getPaste)
	e.POST("/paste", postPaste)

	e.Logger.Fatal(e.Start(":1323"))
}

func index(c echo.Context) error {
	dirEntry, err := os.ReadDir("pastes")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var pastes []string
	for _, entry := range dirEntry {
		pastes = append(pastes, entry.Name())
	}
	return c.Render(http.StatusOK, "index", pastes)
}

func getPaste(c echo.Context) error {
	contentId := c.Param("id")
	content, err := os.ReadFile(fmt.Sprintf("pastes/%s", contentId))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(content))
}

func postPaste(c echo.Context) error {
	content := c.FormValue("content")
	contentId := uuid.New().String()

	if err := os.WriteFile(fmt.Sprintf("pastes/%s", contentId), []byte(content), 0644); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("<a href=\"/paste/%s\">%s<a>", contentId, contentId))
}
