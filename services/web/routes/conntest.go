package routes

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/julienschmidt/httprouter"
)

func HttpConntest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Printf("[Conntest] %s: %v\n", r.Host, r.URL.Path)

	w.Header().Set("X-Organization", "Nintendo")
	w.Write([]byte(`<html><head><title>HTML Page</title></head><body bgcolor="#FFFFFF">This is test.html page</body></html>`))
}

var (
	conntestURL = "conntest.nintendowifi.net"
)

func ConntestRoutes(app *fiber.App) {
	GET(conntestURL, "/", conntest, app)
}

func conntest(c *fiber.Ctx) error {
	c.Set("X-Organization", "Nintendo")
	return c.Render("test", nil)
}
