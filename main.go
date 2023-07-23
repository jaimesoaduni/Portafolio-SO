package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type File struct {
	Name string
	Path string
}

type Pdf File

const file = "public"

func directories(dir string) ([]string, int) {
	var out []string

	entries, err := os.ReadDir(file + dir)
	if err != nil {
		return out, 404
	}

	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}

	return out, 200
}

func find(dir, ext string) []string {
	var out []string

	dir = file + "/" + dir

	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ext {
			out = append(out, dir + "/" + e.Name())
		}
	}

	return out
}

func pdfs(dir string) []Pdf {
	var out []Pdf

	dir = file + "/" + dir

	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".pdf" {
			out = append(out, Pdf{Name: e.Name(), Path: dir + "/" + e.Name()})
		}
	}

	return out
}

func images(dir string) []string {
	return append(find(dir, ".png"), find(dir, ".jpg")...)
}

func otherFiles(dir string) []File {
	var out []File

	dir = file + "/" + dir

	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) != ".pdf" && filepath.Ext(e.Name()) != ".png" && filepath.Ext(e.Name()) != ".jpg" {
			out = append(out, File{Name: e.Name(), Path: dir + "/" + e.Name()})
		}
	}

	return out
}

func main() {
	/* for _, dir := range directories("") {
		fmt.Println(dir)
	} */

	for _, dir := range pdfs(file + "/1. Portada") {
		fmt.Println(dir)
	}

	engine := html.New("templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public/", file)
	app.Static("/css/", "css")
	app.Static("/js/", "js")
	app.Static("/assets/", "assets")

	app.Get("/*", func(c *fiber.Ctx) error {
		pathStr, err := url.QueryUnescape(c.Path())

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		dirs, code := directories(pathStr)
		pathStr = pathStr[1:]

		if code == 404 {
			return c.SendStatus(code)
		}

		var url = c.OriginalURL()

		if len(url) > 1 {
			url += "/"
		}

		return c.Render("index", fiber.Map{
			"url": url,
			"files":  dirs,
			"path":   pathStr,
			"other":  otherFiles(pathStr),
			"images": images(pathStr),
			"pdfs": pdfs(pathStr),
		})
	})

	log.Fatal(app.Listen(":3000"))
}
