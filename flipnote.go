package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func httpFlipnoteRootGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	httpFlipnoteRoot(w, r, ps, false)
}
func httpFlipnoteRootPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	httpFlipnoteRoot(w, r, ps, true)
}
func httpFlipnoteRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params, isPost bool) {
	w.Header()["Content-Type"] = []string{"text/plain"} //Assume plain, except for HTML

	switch ps.ByName("namespace") {
	case "ds":
		switch ps.ByName("region") {
		case "v2-us", "v2-eu", "v2-jp":
			switch ps.ByName("index") {
			case "auth":
				if isPost {
					httpFlipnoteAuthPOST(w, r, ps)
				} else {
					httpFlipnoteAuthGET(w, r, ps)
				}
			case "index.ugo":
				fmt.Println("[Flipnote] UGO stub for index")
			case "en", "es", "jp":
				switch ps.ByName("i1") {
				case "eula.txt":
					w.Write(WriteUTF16String(`Flipnote Hatena has ended its service.
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
						w.Write(WriteUTF16String(`This Flipnote will be deleted from the server.
Deleted material cannot be restored.
Flipnotes that have been downloaded or revised and reposted by other parties will not be affected.`))
						fmt.Println("[Flipnote] Served delete notice")
					case "download.txt":
						w.Write(WriteUTF16String(`Unless specifically regulated by the Terms of Use, please do not use saved Flipnotes on Flipnote Studio, the Flipnote Hatena Website or Flipnote Hatena on the Nintendo DSi for any purposes not legal or not intended for the enjoyment of others.`))
						fmt.Println("[Flipnote] Served download notice")
					case "upload.txt":
						w.Write(WriteUTF16String(`When posting flipnotes you agree to:
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

func httpFlipnoteAuthGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//We have to write to the header map directly so Go doesn't change our casing when using .Set()
	rng := NewUniqueRand(9999999999)
	w.Header()["X-DSi-Auth-Challenge"] = []string{fmt.Sprintf("%d", rng.Int())} //8-10 ASCII characters
	w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
	w.WriteHeader(http.StatusOK) //Necessary when we aren't writing a body
	fmt.Printf("[Auth] Sent auth response: %v\n", w.Header())
}
func httpFlipnoteAuthPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

type UniqueRand struct {
    generated   map[int]bool    //keeps track of
    rng         *rand.Rand      //underlying random number generator
    scope       int             //scope of number to be generated
}

//Generating unique rand less than N
//If N is less or equal to 0, the scope will be unlimited
//If N is greater than 0, it will generate (-scope, +scope)
//If no more unique number can be generated, it will return -1 forwards
func NewUniqueRand(N int) *UniqueRand {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return &UniqueRand{
        generated: map[int]bool{},
        rng:        r1,
        scope:      N,
    }
}

func (u *UniqueRand) Int() int {
    if u.scope > 0 && len(u.generated) >= u.scope {
        return -1
    }
    for {
        var i int
        if u.scope > 0 {
            i = u.rng.Int() % u.scope
        }else{
            i = u.rng.Int()
        }
        if !u.generated[i] {
            u.generated[i] = true
            return i
        }
    }
}