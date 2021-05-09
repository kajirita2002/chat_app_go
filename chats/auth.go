package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// cookieの有無を調べる
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// Userはログインのページにリダイレクトされる
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

// 任意のhttp.Handlerをラップした*authHandlerが生成できる
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// 内部情報を保持する必要がないためhttp.Handlerは実装しなくて良い
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		// エンドポイントが一つ(/chat)なので内容状態は変える必要がない
		// 第三引数では追加の情報をプロバイダから入手したい時に使う
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBiginAuthURLの呼び出し中にエラーが発生しました:", provider, "-", err)
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		// RawQueryの値を解析して認証のプロセスを完了させ認証情報を発行
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした:", provider, "-", err)
		}
		// userの情報を取得する これはjsonデータ
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーを取得できませんでした:", provider, "-", err)
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		userID := fmt.Sprintf("%x", m.Sum(nil))
		// データのnameフィールドの部分をBase64にエンコードする
		// Base64では特殊な文字が入るのを防ぐ(URLやクッキーにセットしたい時に使える)
		authCookieValue := objx.New(map[string]interface{} {
			"userid": userID,
			"name": user.Name(),
			"avatar_url": user.AvatarURL(),
			"email": user.Email(),
		}).MustBase64()
		// authというクッキーにデータを保持する
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		// リダイレクト
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
