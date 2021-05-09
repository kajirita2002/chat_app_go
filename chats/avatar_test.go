package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{
		"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346",
	}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません。")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("Gravatar.GetAvatarURLが%sという誤った値を返しました", url)
	}

}

func TestFileSystemAvatar(t *testing.T) {
	// テスト用のアバターのファイルを生成する
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	// deferの関数は異常終了したとしても必ず実行される
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという誤った値を返しました。", url)
	}
}
