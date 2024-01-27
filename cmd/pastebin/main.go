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

const (
	gitKeepFileName = ".gitkeep"
	pastesDirName   = "pastes"
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
	dirEntry, err := os.ReadDir(pastesDirName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to fetch the pastes!")
	}

	var pastes []string
	for _, entry := range dirEntry {
		if entry.Name() == gitKeepFileName {
			continue
		}
		pastes = append(pastes, entry.Name())
	}
	return c.Render(http.StatusOK, "index", pastes)
}

func getPaste(c echo.Context) error {
	contentId := c.Param("id")
	if contentId == gitKeepFileName {
		return c.String(http.StatusNotFound, "Paste not found!")
	}

	filePath := fmt.Sprintf("%s/%s", pastesDirName, contentId)
	content, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "Paste not found!")
		}
		return c.String(http.StatusInternalServerError, "Unable to fetch the paste!")
	}

	return c.String(http.StatusOK, string(content))
}

func postPaste(c echo.Context) error {
	content := c.FormValue("content")
	contentId := uuid.New().String()

	filePath := fmt.Sprintf("%s/%s", pastesDirName, contentId)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return c.String(http.StatusInternalServerError, "Unable to write the paste!")
	}

	return c.String(http.StatusOK, fmt.Sprintf("<a href=\"/paste/%s\">%s<a>", contentId, contentId))
}
