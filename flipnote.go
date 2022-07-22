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
	fmt.Printf("[HTTP] [Flipnote:Root] %s: %v = %v\n", r.Host, r.URL.Path, r.Header)

	switch ps.ByName("index") {
	case "ds":
		switch ps.ByName("region") {
		case "v2-us", "v2-eu", "v2-jp":
			switch ps.ByName("language") {
			case "auth":
				if isPost {
					httpFlipnoteAuthPOST(w, r, ps)
				} else {
					httpFlipnoteAuthGET(w, r, ps)
				}
			case "en", "es", "jp":
				switch ps.ByName("page") {
				case "eula.txt":
					fmt.Printf("Requested EULA, but none was provided")
				default:
					fmt.Printf("Unknown page: %s", ps.ByName("page"))
				}
			default:
				fmt.Printf("Unknown language: %s", ps.ByName("language"))
			}
		default:
			fmt.Printf("Unknown region: %s", ps.ByName("region"))
		}
	default:
		fmt.Printf("Unknown index: %s", ps.ByName("index"))
	}
}

func httpFlipnoteAuthGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//We have to write to the header map directly so Go doesn't change our casing when using .Set()
	rng := NewUniqueRand(9999999999)
	w.Header()["X-DSi-Auth-Challenge"] = []string{fmt.Sprintf("%d", rng.Int())} //8-10 ASCII characters
	w.Header()["X-DSi-SID"] = []string{"asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"}
	w.WriteHeader(http.StatusOK) //Necessary when we aren't writing a body
	fmt.Printf("[HTTP] [Flipnote:Auth] Sent auth response: %v\n", w.Header())
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
		fmt.Printf("[HTTP] [Flipnote:Auth] Sent auth response: %v\n", w.Header())
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