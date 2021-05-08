package main

import "testing"

// GetAvatarURLメソッドのテスト
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	// 空のclientを実装する
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	// ErrNoAvataraURLが実際にしっかり動くかを検証
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合, AuthAvatar.GetAvatarは" +
			"ErrNoAvatarURLを返すべきです")
	}
	// 値をセットします
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLは" +
		"エラーを返すべきではありません")
	} else {
		// 正しい値が返されるか
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarは正しいURLを返すべきです。")
		}
	}

}