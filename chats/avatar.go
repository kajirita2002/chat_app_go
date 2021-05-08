package main

import "errors"

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