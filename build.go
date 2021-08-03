package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func build() error {
	log.Println("building")

	_ = os.RemoveAll(filepath.Join(".", "build"))

	err := os.MkdirAll(filepath.Join(".", "build", "pygments-cache"), 0755)
	if err != nil {
		return err
	}

	// copy assets
	assetRoot := filepath.Join(".", "assets")
	err = filepath.Walk(assetRoot, func(p string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		dstPath := filepath.Join(".", "build", "assets", p[len(assetRoot):])
		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		src, err := os.Open(p)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// compile pages
	tplRoot := filepath.Join(".", "tpl")
	err = filepath.Walk(tplRoot, func(p string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".gohtml" {
			return nil
		}

		nm := p[len(tplRoot):]
		nm = nm[:len(nm)-len(filepath.Ext(nm))]
		if strings.HasPrefix(nm, string(filepath.Separator)) {
			nm = nm[1:]
		}

		dstPath := filepath.Join(".", "build", nm+".html")
		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		contents, err := compileTemplate(nm, map[string]string{
			"StackVersion": StackVersion,
		})
		if err != nil {
			return err
		}

		_, err = io.WriteString(dst, contents)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func compileCodeBlock(lexer string, src string) (string, error) {
	nm := filepath.Join(".", "build", "pygments-cache", hash([]byte(lexer+"|"+src)))
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

	err = os.WriteFile(nm, codeResult.Bytes(), 0777)
	if err != nil {
		return "", err
	}

	return codeResult.String(), nil
}

func compileTemplate(name string, data interface{}) (string, error) {
	files := []string{}
	err := filepath.Walk("tpl", func(p string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".gohtml") {
			files = append(files, p)
		}
		return err
	})
	if err != nil {
		return "", err
	}

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
