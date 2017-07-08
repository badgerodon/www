package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/badgerodon/rbsa"
	"golang.org/x/net/html"
)

const StackVersion = "v0.3"

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func compileCodeBlock(lexer string, src string) (string, error) {
	nm := filepath.Join("assets", "pygments-cache", hash([]byte(lexer+"|"+src)))
	f, err := os.Open(nm)
	if err == nil {
		bs, err := ioutil.ReadAll(f)
		f.Close()
		if err == nil {
			return string(bs), nil
		}
	}

	src = strings.TrimSpace(src)
	var codeResult bytes.Buffer
	cmd := exec.Command("pygmentize", "-f", "html", "-l", lexer)
	cmd.Stdin = strings.NewReader(src)
	cmd.Stdout = &codeResult
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	ioutil.WriteFile(nm, codeResult.Bytes(), 0777)
	return codeResult.String(), nil
}

func compileTemplate(name string, data interface{}) (string, error) {
	files := []string{}
	filepath.Walk("tpl", func(p string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".gohtml") {
			files = append(files, p)
		}
		return err
	})

	tpl, err := template.New("template").Funcs(template.FuncMap{
		"asset_url": assetURL,
	}).ParseFiles(files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}
	body := buf.String()
	buf.Reset()
	err = tpl.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"Body": template.HTML(body),
	})
	if err != nil {
		return "", err
	}

	var out bytes.Buffer

	inCode := false
	codeLang := "text"
	codeText := ""

	z := html.NewTokenizer(&buf)
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			break loop
		case html.TextToken:
			if inCode {
				codeText += string(z.Text())
			} else {
				out.Write(z.Text())
			}
		case html.StartTagToken:
			tag, hasAttr := z.TagName()
			attrs := make(map[string]string)
			for hasAttr {
				var key, val []byte
				key, val, hasAttr = z.TagAttr()
				attrs[string(key)] = string(val)
			}

			if bytes.Equal(tag, []byte("code")) && attrs["data-lexer"] != "" {
				inCode = true
				codeLang = attrs["data-lexer"]
				codeText = ""
			} else {
				out.WriteString("<")
				out.Write(tag)
				for k, v := range attrs {
					out.WriteString(" ")
					out.WriteString(k)
					out.WriteString("=\"")
					out.WriteString(strings.Replace(v, "\"", "&quot;", -1))
					out.WriteString("\"")
				}
				out.WriteString(">")
			}
		case html.EndTagToken:
			if inCode {
				code, err := compileCodeBlock(codeLang, codeText)
				if err != nil {
					return "", err
				}
				out.WriteString(code)
				inCode = false
			} else {
				out.WriteString(z.Token().String())
			}
		case html.SelfClosingTagToken:
			out.WriteString(z.Token().String())
		case html.CommentToken:
		case html.DoctypeToken:
			out.WriteString(z.Token().String())
		}
	}
	return out.String(), nil
}

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
	fi, err := os.Stat(filepath.Join("assets", name))
	if err != nil {
		return "/assets/0/" + name
	}

	version := fmt.Sprint(fi.ModTime().Unix() + 1)

	return "/assets/" + version + "/" + name
}

func init() {
	log.SetFlags(0)

	// REDIRECTS
	http.Handle("/tools/rbsa/", http.RedirectHandler("/rbsa", 301))

	// API
	http.HandleFunc("/api/rbsa", func(res http.ResponseWriter, req *http.Request) {
		sym := req.FormValue("symbol")
		sym = strings.ToUpper(sym)

		var obj struct {
			Data    map[string]float64 `json:"data"`
			Indices map[string]string  `json:"indices"`
		}

		if sym != "" {
			data, err := rbsa.Analyze(sym)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}

			indices := make(map[string]string)
			for k, _ := range data {
				indices[k] = rbsa.DEFAULT_INDICES[k]
			}

			obj.Data = data
			obj.Indices = indices
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(obj)
	})

	// HTTP
	staticHandler := http.FileServer(http.Dir("./assets"))
	http.HandleFunc("/assets/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Cache-Control", "max-age=31556926")
		parts := strings.Split(req.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(res, "not found", 404)
			return
		}
		req.URL.Path = strings.Join(parts[3:], "/")
		staticHandler.ServeHTTP(res, req)
	})
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on 127.0.0.1:%v\n", port)
	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
