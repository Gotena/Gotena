package web

import (
	"github.com/JoshuaDoes/gotena/services/web/routes"
	"github.com/gofiber/fiber/v2"
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

	initRoutes(app)

	return app
}

func initRoutes(app *fiber.App) {
	routes.ConntestRoutes(app)
	routes.FlipnoteRoutes(app)
}
