package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const StackVersion = "v0.3"

func serveTemplate(res http.ResponseWriter, name string, data interface{}) {
	contents, err := compileTemplate(name, data)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	io.WriteString(res, contents)
}

func assetURL(name string) string {
	return "/assets/" + name
}

func init() {
	log.SetFlags(0)

	// HTTP
	http.HandleFunc("/stack", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		serveTemplate(res, "stack/index", map[string]string{
			"StackVersion": StackVersion,
		})
	})
	http.HandleFunc("/license", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		serveTemplate(res, "license", nil)
	})
	http.HandleFunc("/rbsa", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		serveTemplate(res, "rbsa/index", nil)
	})
	http.HandleFunc("/socketmaster", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		serveTemplate(res, "socketmaster/index", nil)
	})
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" || strings.ToLower(req.URL.Path) == "/index.htm" {
			res.Header().Set("Content-Type", "text/html")
			serveTemplate(res, "index", nil)
		} else {
			http.Error(res, "Page Not Found", 404)
		}
	})
}

func main() {
	log.SetFlags(0)
	err := build()
	if err != nil {
		log.Fatalln(err)
	}
}
