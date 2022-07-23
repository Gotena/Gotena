package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func httpConntest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Printf("[Conntest] %s: %v\n", r.Host, r.URL.Path)
	
	w.Header().Set("X-Organization", "Nintendo")
	w.Write([]byte(`<html><head><title>HTML Page</title></head><body bgcolor="#FFFFFF">This is test.html page</body></html>`))
}