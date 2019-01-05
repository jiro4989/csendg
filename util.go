package main

import (
	"os"
)

// HomeDir はホームディレクトリを取得する。
// HOMEがあればHOMEを返し、USERPROFILEがあればUSERPROFILEを返す。
// 両方なければエラーでアプリが死ぬ
func HomeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			panic("環境変数が定義されてません。環境変数を確認してください。")
		}
		return home
	}
	return home
}

// Exists はファイルの有無を返す。
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
