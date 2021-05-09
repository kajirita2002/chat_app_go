package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarインターフェースがアバターのURLを返すことができない場合に発生するエラーです
// 一回だけ生成される
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターにURLを返します。
	// 問題が発生した場合はエラーを返します。特に、URLの取得できなかった場合には
	// ErrNoAvatarURLを返します。
	GetAvatarURL(c *client) (string, error)
}

// 空の構造体として定義
// フィールドがないのでメソッドはレシーバーを参照する必要がない
type AuthAvatar struct{}

// AuthAvatarを利用してチャットルームを新規作成する
// ヘルパーとして機能する
var UseAuthAvatar AuthAvatar

// オブジェクトにフィールドがないためレシーバーを参照する必要がない
// Userデータからurlを取得する
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	// データがあるかどうか
	if url, ok := c.userData["avatar_url"]; ok {
		// 文字列かどうか
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

// GravatarAvatar型の利用を容易にするため
var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			// 全てのファイルのリストを取得する
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					// IsDirはディレクトリかどうかを確認し、処理対象からディレクトリを除外する
					if file.IsDir() {
						continue
					}
					// *はその他の任意の文字列にマッチしていれば良い
					if match, _ := filepath.Match(useridStr+"*", file.Name());
						match {
							return "/avatars/" + file.Name(), nil
						}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
