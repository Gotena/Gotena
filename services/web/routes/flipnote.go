package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Gotena/Gotena/tools"
	"github.com/Gotena/Gotena/utils"
	"github.com/RinLovesYou/ppmlib-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

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

// func HttpFlipnoteRootGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	HttpFlipnoteRoot(w, r, ps, false)
// }
// func HttpFlipnoteRootPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	HttpFlipnoteRoot(w, r, ps, true)
// }
// func HttpFlipnoteRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params, isPost bool) {
// 	w.Header()["Content-Type"] = []string{"text/plain"} //Assume plain, except for HTML

// 	switch ps.ByName("namespace") {
// 	case "ds":
// 		switch ps.ByName("region") {
// 		case "v2-us", "v2-eu", "v2-jp":
// 			switch ps.ByName("index") {
// 			case "auth":
// 				if isPost {
// 					HttpFlipnoteAuthPOST(w, r, ps)
// 				} else {
// 					HttpFlipnoteAuthGET(w, r, ps)
// 				}
// 			case "index.ugo":
// 				fmt.Println("[Flipnote] UGO stub for index")
// 			case "en", "es", "jp":
// 				switch ps.ByName("i1") {
// 				case "eula.txt":
// 					w.Write(utils.WriteUTF16String(`Flipnote Hatena has ended its service.
// This server is written and hosted by
// JoshuaDoes and RinLovesYou. Our source:
// https://github.com/Gotena/Gotena

// This is still in very early stages, so
// please be gentle with my server. Or don't,
// we want to validate good data!`))
// 					fmt.Println("[Flipnote] Served EULA")
// 				case "confirm":
// 					switch ps.ByName("i2") {
// 					case "delete.txt":
// 						w.Write(utils.WriteUTF16String(`This Flipnote will be deleted from the server.
// Deleted material cannot be restored.
// Flipnotes that have been downloaded or revised and reposted by other parties will not be affected.`))
// 						fmt.Println("[Flipnote] Served delete notice")
// 					case "download.txt":
// 						w.Write(utils.WriteUTF16String(`Unless specifically regulated by the Terms of Use, please do not use saved Flipnotes on Flipnote Studio, the Flipnote Hatena Website or Flipnote Hatena on the Nintendo DSi for any purposes not legal or not intended for the enjoyment of others.`))
// 						fmt.Println("[Flipnote] Served download notice")
// 					case "upload.txt":
// 						w.Write(utils.WriteUTF16String(`When posting flipnotes you agree to:
// ・ The flipnote becoming freely available to others on the internet
// ・ Others will be able to save your flipnote and convert it to video
// ・ Unlocked flipnotes can be altered and/or reposted by others easily
// ・ No locked flipnote is truly locked for a determined theft, report stolen content`))
// 						fmt.Println("[Flipnote] Served upload notice")
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func HttpFlipnoteAuthGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	//We have to write to the header map directly so Go doesn't change our casing when using .Set()
// 	rng := utils.NewUniqueRand(9999999999)
// 	w.Header()["X-DSi-Auth-Challenge"] = []string{fmt.Sprintf("%d", rng.Int())} //8-10 ASCII characters
// 	w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
// 	w.WriteHeader(http.StatusOK) //Necessary when we aren't writing a body
// 	fmt.Printf("[Auth] Sent auth response: %v\n", w.Header())
// }
// func HttpFlipnoteAuthPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	if r.Header.Get("X-DSi-ID") == "" {
// 		fmt.Printf("[HTTP] [Flipnote:Auth] Invalid auth attempt\n")
// 		w.WriteHeader(http.StatusForbidden)
// 	} else {
// 		w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
// 		w.Header()["X-DSi-New-Notices"] = []string{"1"}
// 		w.Header()["X-DSi-Unread-Notices"] = []string{"1"}
// 		w.WriteHeader(http.StatusOK)
// 		fmt.Printf("[Auth] Sent auth response: %v\n", w.Header())
// 	}
// }

var (
	flipnoteURL = "flipnote.hatena.com"
)

func FlipnoteRoutes(app *fiber.App) {
	GET(flipnoteURL, "/ds/:region/auth", flipnoteAuthGet, app)
	POST(flipnoteURL, "/ds/:region/auth", flipnoteAuthPost, app)

	GET(flipnoteURL, "/ds/:region/index.ugo", flipnoteIndexGet, app)
	GET(flipnoteURL, "/ds/:region/frontpage/frontpage.ugo", flipnoteFrontPageGet, app)
	GET(flipnoteURL, "/ds/:region/frontpage/hotmovies.ugo", flipnoteHotMoviesGet, app)

	GET(flipnoteURL, "/ds/:region/movie/:id/:file", flipnoteMovieGet, app)

	GET(flipnoteURL, "/ds/:region/:language/eula.txt", flipnoteEulaGet, app)
	GET(flipnoteURL, "/ds/:region/:language/confirm/:file", flipnoteConfirmGet, app)
	GET(flipnoteURL, "/ds/:region/help/info.htm", flipnoteInfoGet, app)
	GET(flipnoteURL, "css/ds/:file", flipnoteCssGet, app)
}

func flipnoteCssGet(c *fiber.Ctx) error {
	file := c.Params("file")
	if file == "" {
		c.Status(http.StatusNotFound)
		return c.SendStatus(http.StatusNotFound)
	}

	return c.SendFile(fmt.Sprintf("services/web/routes/res/css/%s", file))
}

func flipnoteAuthGet(c *fiber.Ctx) error {
	region := c.Params("region")
	if region != "v2-us" && region != "v2-eu" && region != "v2-jp" {
		fmt.Println("[Flipnote] Invalid region:", region)
		return c.SendStatus(http.StatusBadRequest)
	}

	sess, err := c.Locals("sessions").(*session.Store).Get(c)
	if err != nil {
		fmt.Println("[Fiber] Unable to retrieve session:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	authChallenge := fmt.Sprintf("%d", utils.NewUniqueRand(9999999999).Int()) //~10 ASCII characters
	sess.Set("auth", authChallenge)

	if err := sess.Save(); err != nil {
		fmt.Printf("[Fiber] Error saving session for SID %s: %v\n", sess.ID(), err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	fmt.Println("[Flipnote] Starting session with auth challenge:", authChallenge)
	c.Context().SetContentType("text/plain")
	c.Response().Header.Set("X-DSi-Auth-Challenge", authChallenge)
	c.Response().Header.SetStatusCode(http.StatusOK) //Necessary when we aren't writing a body

	return nil
}

func flipnoteAuthPost(c *fiber.Ctx) error {
	region := c.Params("region")
	if region != "v2-us" && region != "v2-eu" && region != "v2-jp" {
		fmt.Println("[Flipnote] Invalid region:", region)
		return c.SendStatus(http.StatusBadRequest)
	}

	sess, err := c.Locals("sessions").(*session.Store).Get(c)
	if err != nil {
		fmt.Println("[Fiber] Unable to retrieve session:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	fsid := c.Get("X-DSi-ID")
	if fsid == "" {
		fmt.Println("[Flipnote] Invalid auth attempt, missing FSID")
	}

	authResp := c.Get("X-DSi-Auth-Response")
	if authResp == "" {
		fmt.Println("[Flipnote] Invalid auth attempt, missing auth challenge response")
		return c.SendStatus(http.StatusForbidden)
	}

	authChallenge := sess.Get("auth")
	if authChallenge == nil {
		fmt.Println("[Fiber] Error performing auth challenge: no auth attempt in session")
		return c.SendStatus(http.StatusInternalServerError)
	}

	authCheck, err := tools.AuthChallenge(fsid, authChallenge.(string))
	if err != nil {
		fmt.Println("[Flipnote] Error performing auth challenge:", err)
		return c.SendStatus(http.StatusForbidden)
	}

	if authCheck != authResp {
		fmt.Printf("[Flipnote] Invalid auth attempt, check (%s) does not match response (%s)\n", authCheck, authResp)
		return c.SendStatus(http.StatusForbidden)
	}

	sess.Set("authResp", authResp)
	if err := sess.Save(); err != nil {
		fmt.Printf("[Fiber] Error saving session for SID %s: %v\n", sess.ID(), err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Response().Header.Set("X-DSi-SID", sess.ID())
	c.Response().Header.Set("X-DSi-New-Notices", "1")
	c.Response().Header.Set("X-DSi-Unread-Notices", "1")
	c.Response().SetStatusCode(http.StatusOK)

	return nil
}

func flipnoteIndexGet(c *fiber.Ctx) error {
	ugoBytes, err := tools.PackUgoJson("services/web/routes/res/ugo/index.ugo.json")
	if err != nil {
		fmt.Println("[Flipnote] Error reading UGO for index:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.Send(ugoBytes)
}

func flipnoteFrontPageGet(c *fiber.Ctx) error {
	ugoBytes, err := tools.PackUgoJson("services/web/routes/res/ugo/frontpage.ugo.json")
	if err != nil {
		fmt.Println("[Flipnote] Error reading UGO for front page:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.Send(ugoBytes)
}

func flipnoteEulaGet(c *fiber.Ctx) error {
	lang := c.Params("language")
	if lang != "en" && lang != "jp" && lang != "es" {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.Send(utils.WriteUTF16String(eula))
}

func flipnoteConfirmGet(c *fiber.Ctx) error {
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

func flipnoteInfoGet(c *fiber.Ctx) error {
	c.Response().Header.Set("content-type", "text/html")

	file, err := os.ReadFile("services/web/views/info.htm")
	if err != nil {
		fmt.Println("[Flipnote] Error reading info view:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Send(file)
}

func flipnoteHotMoviesGet(c *fiber.Ctx) error {
	ugoBytes, err := tools.PackUgoJson("services/web/routes/res/ugo/hotmovies.ugo.json")
	if err != nil {
		fmt.Println("[Flipnote] Error reading UGO for hot movies:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.Send(ugoBytes)
}

func flipnoteMovieGet(c *fiber.Ctx) error {
	c.Response().Header.Set("content-type", "text/html")

	file, err := os.ReadFile("services/web/views/detailsTemplate.htm")
	if err != nil {
		fmt.Println("[help] Error reading info file:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	entrySeparatorFile, err := os.ReadFile("services/web/views/pageEntrySeparator.htm")
	if err != nil {
		fmt.Println("[help] Error reading info file:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	entrySeparator := string(entrySeparatorFile)

	entryTemplateFile, err := os.ReadFile("services/web/views/pageEntryTemplate.htm")
	if err != nil {
		fmt.Println("[help] Error reading info file:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	entryTemplate := string(entryTemplateFile)

	creator := c.Params("id")
	filename := c.Params("file")

	fmt.Println("[Flipnote] Movie request:", creator, filename)

	if strings.HasSuffix(filename, ".info") {
		c.Response().Header.Set("content-type", "text/plain")
		fmt.Println("[Flipnote] Movie info request:", creator, filename)
		return c.SendString("0\n0\n")
	}

	if strings.HasSuffix(filename, ".ppm") {
		c.Response().Header.Set("content-type", "text/plain")
		file, err := os.ReadFile("services/web/routes/res/" + filename)
		if err != nil {
			fmt.Println("[Flipnote] Error reading PPM to serve:", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Send(file)
	}

	if strings.HasSuffix(filename, ".htm") {
		filename = strings.TrimSuffix(filename, ".htm")

		ppmFile, err := os.ReadFile("services/web/routes/res/" + filename + ".ppm")
		if err != nil {
			fmt.Println("[Flipnote] Error reading PPM to parse:", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		ppm, err := ppmlib.Parse(ppmFile)
		if err != nil {
			fmt.Println("[Flipnote] Error parsing PPM:", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		entries := make([]string, 0)

		//Creator Username
		creatorContent := fmt.Sprintf(`<a href="http://flipnote.hatena.com/ds/v2-xx/%s/profile.htm?t=260&pm=80\">%s</a>`, creator, ppm.CurrentAuthor.Name)
		usernameEntry := entryTemplate
		usernameEntry = strings.ReplaceAll(usernameEntry, "{{Name}}", "Creator")
		usernameEntry = strings.ReplaceAll(usernameEntry, "{{Content}}", creatorContent)
		entries = append(entries, usernameEntry)

		//Stars
		stars := "0"
		starsContent := fmt.Sprintf(`<a href="http://flipnote.hatena.com/ds/v2-xx/movie/%s/%s.htm?mode=stardetail"><span class="star0c">★</span> <span class="star0">%s</span></a>`, creator, filename, stars)
		starsContent += fmt.Sprintf(`<br/><a href="http://flipnote.hatena.com/ds/v2-xx/movie/%s/%s.htm?mode=stardetail"><span class="star1c">★</span> <span class="star1">%s</span></a>`, creator, filename, stars)
		starsContent += fmt.Sprintf(`<br/><a href="http://flipnote.hatena.com/ds/v2-xx/movie/%s/%s.htm?mode=stardetail"><span class="star2c">★</span> <span class="star2">%s</span></a>`, creator, filename, stars)
		starsContent += fmt.Sprintf(`<br/><a href="http://flipnote.hatena.com/ds/v2-xx/movie/%s/%s.htm?mode=stardetail"><span class="star3c">★</span> <span class="star3">%s</span></a>`, creator, filename, stars)
		starsContent += fmt.Sprintf(`<br/><a href="http://flipnote.hatena.com/ds/v2-xx/movie/%s/%s.htm?mode=stardetail"><span class="star4c">★</span> <span class="star4">%s</span></a>`, creator, filename, stars)
		starsEntry := entryTemplate
		starsEntry = strings.ReplaceAll(starsEntry, "{{Name}}", "Stars")
		starsEntry = strings.ReplaceAll(starsEntry, "{{Content}}", starsContent)
		entries = append(entries, starsEntry)

		//Views
		viewsContent := "0"
		viewsEntry := entryTemplate
		viewsEntry = strings.ReplaceAll(viewsEntry, "{{Name}}", "Views")
		viewsEntry = strings.ReplaceAll(viewsEntry, "{{Content}}", viewsContent)
		entries = append(entries, viewsEntry)

		//Downloads
		downloadsContent := "0"
		downloadsEntry := entryTemplate
		downloadsEntry = strings.ReplaceAll(downloadsEntry, "{{Name}}", "Downloads")
		downloadsEntry = strings.ReplaceAll(downloadsEntry, "{{Content}}", downloadsContent)
		entries = append(entries, downloadsEntry)

		//Channel
		channelContent := fmt.Sprintf(`<a href="http://flipnote.hatena.com/ds/v2-xx/ch/%s.uls">%s</a>`, "idk", "idk")
		channelEntry := entryTemplate
		channelEntry = strings.ReplaceAll(channelEntry, "{{Name}}", "Channel")
		channelEntry = strings.ReplaceAll(channelEntry, "{{Content}}", channelContent)
		entries = append(entries, channelEntry)

		fileString := string(file)
		fileString = strings.ReplaceAll(fileString, "{{CreatorID}}", creator)
		fileString = strings.ReplaceAll(fileString, "{{Filename}}", filename)
		fileString = strings.ReplaceAll(fileString, "{{CommentCount}}", "0")
		fileString = strings.ReplaceAll(fileString, "{{PageEntries}}", strings.Join(entries, entrySeparator))
		fileString = strings.ReplaceAll(fileString, "{{Username}}", ppm.CurrentAuthor.Name)

		return c.Send([]byte(fileString))

	}

	return c.SendStatus(http.StatusNotFound)
}
