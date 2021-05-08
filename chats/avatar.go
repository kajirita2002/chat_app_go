package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			// MD5を使って ハッシュ値を算出する
			m := md5.New()
			// hash.hashはio.Writerインターフェースを実装しているためio.WriteStrngメソッドで文字列を与えられる
			// 大文字を小文字に変換したものをmに書き込む
			io.WriteString(m, strings.ToLower(emailStr))
			// sumでその時点までに書き込まれた文字列を使ってハッシュ値が計算される
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}