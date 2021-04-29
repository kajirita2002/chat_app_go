package main

import "net/http"

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
