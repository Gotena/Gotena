package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func httpConntest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Printf("[HTTP] [Conntest] %s: %v\n", r.Host, r.URL.Path)
	
	w.Header().Set("Server", "BigIP")
	w.Header().Set("X-Organization", "Nintendo")
	w.Write([]byte(`
            <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
            <html>
            <head>
            <title>HTML Page</title>
            </head>
            <body bgcolor="#FFFFFF">
            This is test.html page
            </body>
            </html>
          `))
}