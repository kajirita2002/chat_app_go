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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// 現在アクティブなAvatarの実装
var avatars Avatar = UseFileSystemAvatar

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// dataをtemplateに渡す
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "The host of the application")

func main() {
	flag.Parse()
	gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		facebook.New("2891035937832238", "dc17fa42bf71f058cd1c8246d3c701e4", "http://localhost:8080/auth/callback/facebook"),
		github.New("4c86aced996c4c012f75", "258d77a40adf1367faec823baa08f0cd58f04228", "http://localhost:8080/auth/callback/github"),
		google.New("441470423311-l3va9r1oo3mh0uqtpaehnclpqkb4m90n.apps.googleusercontent.com", "cd6JZE8tagHAT_2NycTfGl8-", "http://localhost:8080/auth/callback/google"),
	)
	// AuthAvatarのインスタンスを作成していないためメモリ使用量が増えることはない
	r := newRoom()
	// os.Stoutでターミナルに出力が行われる
	r.tracer = trace.New(os.Stdout)
	// ます*authHandlerのServeHTTPメソッド⇨認証成功⇨*templateHandlerのServeHTTPメソッドが実行
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	// http.Handlerを実装していないハンドラはHandleFunc関数を使う
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: "",
			Path:  "/",
			// クッキーは即座に削除される
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	http.Handle("/avatars/",
		// http.StriPrefix型はhttp.Handlerを受け取ってパスを変更する(接頭辞の部分を削除)
		// /avatars/avatarsになってしまうため
		http.StripPrefix("/avatars/",
			// 静的なファイルの提供やファイル一覧の作成、404エラーの作成などの機能を備えている
			// http.Dirは公開しようとしているフォルダーを指定するために利用される
			http.FileServer(http.Dir("./avatars"))))

	go r.run()

	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
