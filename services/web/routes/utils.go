package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	getRoutes = make(map[string]map[string]func(c *fiber.Ctx) error)
	putRoutes = make(map[string]map[string]func(c *fiber.Ctx) error)
)

func GET(host string, path string, handler func(c *fiber.Ctx) error, app *fiber.App) {
	if val, ok := getRoutes[host]; !ok || val == nil {
		getRoutes[host] = make(map[string]func(c *fiber.Ctx) error)
		getRoutes[host][path] = handler
	} else {
		getRoutes[host][path] = handler
	}

	app.Get(path, getMiddleware, func(c *fiber.Ctx) error {
		basePath := c.BaseURL()
		basePath = strings.ReplaceAll(basePath, "http://", "")
		basePath = strings.ReplaceAll(basePath, "https://", "")
		basePath = strings.Split(basePath, "/")[0]

		if val, ok := getRoutes[basePath]; ok && val != nil {
			if handler, ok := val[path]; ok {
				return handler(c)
			}
		}

		return c.SendStatus(http.StatusNotFound)
	})

}

func POST(host string, path string, handler func(c *fiber.Ctx) error, app *fiber.App) {
	if val, ok := putRoutes[host]; !ok || val == nil {
		putRoutes[host] = make(map[string]func(c *fiber.Ctx) error)
		putRoutes[host][path] = handler
	} else {
		putRoutes[host][path] = handler
	}

	app.Post(path, postMiddleware, func(c *fiber.Ctx) error {
		basePath := c.BaseURL()
		basePath = strings.ReplaceAll(basePath, "http://", "")
		basePath = strings.ReplaceAll(basePath, "https://", "")
		basePath = strings.Split(basePath, "/")[0]

		if val, ok := putRoutes[basePath]; ok && val != nil {
			if handler, ok := val[path]; ok {
				return handler(c)
			}
		}

		return c.SendStatus(http.StatusNotFound)
	})

}

func getMiddleware(c *fiber.Ctx) error {

	fmt.Printf("[HTTP] [GET] [%s]\nHeader: %s", c.BaseURL(), string(c.Request().Header.Header()))

	return c.Next()
}

func postMiddleware(c *fiber.Ctx) error {

	fmt.Printf("[HTTP] [POST] [%s]\nHeader: %s", c.BaseURL(), string(c.Request().Header.Header()))

	return c.Next()
}
