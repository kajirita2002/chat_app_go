package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	// 隠しフィールドとして送信されたユーザーIDの値を読み取る
	userId := r.FormValue("userid")
	// アップロードされたバイト列を読み込むためのio.Reader型の値を取得する
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	// io.Readerか全てのバイト列を読み込みます
	// クライアントから送られたバイト列をdataとして受け取る
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	// useridの値をもとに保存先のファイル名の文字列を生成する
	filename := filepath.Join("avatars", userId+filepath.Ext(header.Filename))
	// avatarsフォルダーに新規ファイルを作成してデータを保存する
	// 0777は全てのユーザーに対してこのファイルへの全てのアクセス権を与えるというもの
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "成功")
}