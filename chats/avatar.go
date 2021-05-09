package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
)

// ErrNoAvatarインターフェースがアバターのURLを返すことができない場合に発生するエラーです
// 一回だけ生成される
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターにURLを返します。
	// 問題が発生した場合はエラーを返します。特に、URLの取得できなかった場合には
	// ErrNoAvatarURLを返します。
	// GetAvatarURLメソッドが受け取れる型をより柔軟にする
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// 空の構造体として定義
// フィールドがないのでメソッドはレシーバーを参照する必要がない
type AuthAvatar struct{}

// AuthAvatarを利用してチャットルームを新規作成する
// ヘルパーとして機能する
var UseAuthAvatar AuthAvatar

// オブジェクトにフィールドがないためレシーバーを参照する必要がない
// Userデータからurlを取得する
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

// GravatarAvatar型の利用を容易にするため
var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	// 全てのファイルのリストを取得する
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			// IsDirはディレクトリかどうかを確認し、処理対象からディレクトリを除外する
			if file.IsDir() {
				continue
			}
			log.Println(u.UniqueID())
			log.Println(file.Name())
			// *はその他の任意の文字列にマッチしていれば良い
			// clientのuserDataフィールドにアクセスする代わりにChatUserインターフェースのUniqueIDメソッドが呼び出されている
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
