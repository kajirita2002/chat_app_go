package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

// GetAvatarURLメソッドのテスト
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	// AvatarURLがエラーを返す
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	// ErrNoAvataraURLが実際にしっかり動くかを検証
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合, AuthAvatar.GetAvatarは" +
			"ErrNoAvatarURLを返すべきです")
	}
	// 値をセットします
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	// AvatarURLがtestUrlを返す
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
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
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません。")
	}
	if url != "//www.gravatar.com/avatar/abc" {
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
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという誤った値を返しました。", url)
	}
}
