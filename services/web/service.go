package web

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Gotena/Gotena/services/web/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html"
)

// func Init() *HTTP {
// 	return NewHTTP().
// 		GET("conntest.nintendowifi.net", "/", routes.HttpConntest).
// 		GET("flipnote.hatena.com", "/:namespace", routes.HttpFlipnoteRootGET).
// 		GET("flipnote.hatena.com", "/:namespace/:region", routes.HttpFlipnoteRootGET).
// 		GET("flipnote.hatena.com", "/:namespace/:region/:index", routes.HttpFlipnoteRootGET).
// 		GET("flipnote.hatena.com", "/:namespace/:region/:index/:i1", routes.HttpFlipnoteRootGET).
// 		GET("flipnote.hatena.com", "/:namespace/:region/:index/:i1/:i2", routes.HttpFlipnoteRootGET).
// 		POST("flipnote.hatena.com", "/:namespace", routes.HttpFlipnoteRootPOST).
// 		POST("flipnote.hatena.com", "/:namespace/:region", routes.HttpFlipnoteRootPOST).
// 		POST("flipnote.hatena.com", "/:namespace/:region/:index", routes.HttpFlipnoteRootPOST).
// 		POST("flipnote.hatena.com", "/:namespace/:region/:index/:i1", routes.HttpFlipnoteRootPOST).
// 		POST("flipnote.hatena.com", "/:namespace/:region/:index/:i1/:i2", routes.HttpFlipnoteRootPOST)
// }

var TemplateEngine *html.Engine

func InitFiber() *fiber.App {
	TemplateEngine = html.New("./services/web/views", ".html")
	TemplateEngine.Delims("{{", "}}")

	app := fiber.New(fiber.Config{
		Views:                    TemplateEngine,
		DisableKeepalive:         false,
		DisableHeaderNormalizing: true,
	})

	app.Use(func (c *fiber.Ctx) error {
		c.Locals("sessions", session.New(session.Config{
			KeyGenerator: func() string {
				//Note - NOT RFC4122 compliant
				b := make([]byte, 20)
				_, err := rand.Read(b)
				if err != nil {
					return ""
				}
				return fmt.Sprintf("%X", b)
			},
			KeyLookup: "header:X-DSi-SID",
			Storage: sqlite3.New(sqlite3.Config{
				Database: "./gotena.sqlite3",
				Table: "fiber_sessions",
				Reset: false,
				GCInterval: 10 * time.Second,
				MaxOpenConns: 100,
				MaxIdleConns: 100,
				ConnMaxLifetime: 1 * time.Second,
			}),
		}))
		return c.Next()
	})

	initRoutes(app)

	return app
}

func initRoutes(app *fiber.App) {
	routes.ConntestRoutes(app)
	routes.FlipnoteRoutes(app)
}
