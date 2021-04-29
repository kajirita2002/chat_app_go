package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/kaji2002/chat_app/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

var host = flag.String("host", ":8080", "The host of the application")

func main() {
	flag.Parse()
	r := newRoom()
	// os.Stoutでターミナルに出力が行われる
	r.tracer = trace.New(os.Stdout)
	// ます*authHandlerのServeHTTPメソッド⇨認証成功⇨*templateHandlerのServeHTTPメソッドが実行
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
