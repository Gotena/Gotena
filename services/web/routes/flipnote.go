package routes

import (
	"fmt"
	"net/http"

	"github.com/JoshuaDoes/gotena/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/julienschmidt/httprouter"

	_ "embed"
)

var (
	//go:embed res/eula.txt
	eula string

	//go:embed res/delete.txt
	deleteNotice string

	//go:embed res/download.txt
	downloadNotice string

	//go:embed res/upload.txt
	uploadNotice string
)

func HttpFlipnoteRootGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	HttpFlipnoteRoot(w, r, ps, false)
}
func HttpFlipnoteRootPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	HttpFlipnoteRoot(w, r, ps, true)
}
func HttpFlipnoteRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params, isPost bool) {
	w.Header()["Content-Type"] = []string{"text/plain"} //Assume plain, except for HTML

	switch ps.ByName("namespace") {
	case "ds":
		switch ps.ByName("region") {
		case "v2-us", "v2-eu", "v2-jp":
			switch ps.ByName("index") {
			case "auth":
				if isPost {
					HttpFlipnoteAuthPOST(w, r, ps)
				} else {
					HttpFlipnoteAuthGET(w, r, ps)
				}
			case "index.ugo":
				fmt.Println("[Flipnote] UGO stub for index")
			case "en", "es", "jp":
				switch ps.ByName("i1") {
				case "eula.txt":
					w.Write(utils.WriteUTF16String(`Flipnote Hatena has ended its service.
This server is written and hosted by
JoshuaDoes and RinLovesYou. Our source:
https://github.com/Gotena/Gotena

This is still in very early stages, so
please be gentle with my server. Or don't,
we want to validate good data!`))
					fmt.Println("[Flipnote] Served EULA")
				case "confirm":
					switch ps.ByName("i2") {
					case "delete.txt":
						w.Write(utils.WriteUTF16String(`This Flipnote will be deleted from the server.
Deleted material cannot be restored.
Flipnotes that have been downloaded or revised and reposted by other parties will not be affected.`))
						fmt.Println("[Flipnote] Served delete notice")
					case "download.txt":
						w.Write(utils.WriteUTF16String(`Unless specifically regulated by the Terms of Use, please do not use saved Flipnotes on Flipnote Studio, the Flipnote Hatena Website or Flipnote Hatena on the Nintendo DSi for any purposes not legal or not intended for the enjoyment of others.`))
						fmt.Println("[Flipnote] Served download notice")
					case "upload.txt":
						w.Write(utils.WriteUTF16String(`When posting flipnotes you agree to:
・ The flipnote becoming freely available to others on the internet
・ Others will be able to save your flipnote and convert it to video
・ Unlocked flipnotes can be altered and/or reposted by others easily
・ No locked flipnote is truly locked for a determined theft, report stolen content`))
						fmt.Println("[Flipnote] Served upload notice")
					}
				}
			}
		}
	}
}

func HttpFlipnoteAuthGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//We have to write to the header map directly so Go doesn't change our casing when using .Set()
	rng := utils.NewUniqueRand(9999999999)
	w.Header()["X-DSi-Auth-Challenge"] = []string{fmt.Sprintf("%d", rng.Int())} //8-10 ASCII characters
	w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
	w.WriteHeader(http.StatusOK) //Necessary when we aren't writing a body
	fmt.Printf("[Auth] Sent auth response: %v\n", w.Header())
}
func HttpFlipnoteAuthPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Header.Get("X-DSi-ID") == "" {
		fmt.Printf("[HTTP] [Flipnote:Auth] Invalid auth attempt\n")
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
		w.Header()["X-DSi-New-Notices"] = []string{"1"}
		w.Header()["X-DSi-Unread-Notices"] = []string{"1"}
		w.WriteHeader(http.StatusOK)
		fmt.Printf("[Auth] Sent auth response: %v\n", w.Header())
	}
}

var (
	flipnoteURL = "flipnote.hatena.com"
)

func FlipnoteRoutes(app *fiber.App) {
	GET(flipnoteURL, "/ds/:region/auth", flipnoteAuthGet, app)
	POST(flipnoteURL, "/ds/:region/auth", flipnoteAuthPost, app)

	GET(flipnoteURL, "/ds/:region/index.ugo", flipnoteIndexGet, app)

	GET(flipnoteURL, "/ds/:region/:language/eula.txt", flipnoteEulaGet, app)
	GET(flipnoteURL, "/ds/:region/:language/confirm/:file", flipnoteConfirmGet, app)
}

func flipnoteAuthGet(c *fiber.Ctx) error {

	region := c.Params("region")
	if region != "v2-us" && region != "v2-eu" && region != "v2-jp" {
		fmt.Println("[Flipnote] Invalid region:", region)
		return c.SendStatus(http.StatusBadRequest)
	}

	rng := utils.NewUniqueRand(9999999999)

	c.Context().SetContentType("text/plain")
	c.Response().Header.Set("X-DSi-Auth-Challenge", fmt.Sprintf("%d", rng.Int())) //8-10 ASCII characters
	c.Response().Header.Set("X-DSi-SID", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf")
	c.Response().Header.SetStatusCode(http.StatusOK) //Necessary when we aren't writing a body

	return nil
}

func flipnoteAuthPost(c *fiber.Ctx) error {
	region := c.Params("region")
	if region != "v2-us" && region != "v2-eu" && region != "v2-jp" {
		fmt.Println("[Flipnote] Invalid region:", region)
		return c.SendStatus(http.StatusBadRequest)
	}

	id := c.Get("X-DSi-ID")
	if id == "" {
		fmt.Printf("[HTTP] [Flipnote:Auth] Invalid auth attempt\n")
		return c.SendStatus(http.StatusForbidden)
	}

	c.Response().Header.Set("X-DSi-SID", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf")
	c.Response().Header.Set("X-DSi-New-Notices", "1")
	c.Response().Header.Set("X-DSi-Unread-Notices", "1")
	c.Response().SetStatusCode(http.StatusOK)

	return nil
}

func flipnoteIndexGet(c *fiber.Ctx) error {
	fmt.Println("[index] ugo stub")

	return c.SendStatus(http.StatusOK)
}

func flipnoteEulaGet(c *fiber.Ctx) error {
	fmt.Println("[eula] ugo stub")

	lang := c.Params("language")
	if lang != "en" && lang != "jp" && lang != "es" {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.Send(utils.WriteUTF16String(eula))
}

func flipnoteConfirmGet(c *fiber.Ctx) error {
	fmt.Println("[confirm] ugo stub")

	file := c.Params("file")
	lang := c.Params("language")

	if lang != "en" && lang != "jp" && lang != "es" {
		return c.SendStatus(http.StatusBadRequest)
	}

	switch file {
	case "delete.txt":
		return c.Send(utils.WriteUTF16String(deleteNotice))
	case "download.txt":
		return c.Send(utils.WriteUTF16String(downloadNotice))
	case "upload.txt":
		return c.Send(utils.WriteUTF16String(uploadNotice))
	}

	return c.SendStatus(http.StatusNotFound)
}
