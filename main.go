package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

type AppConfig struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

var (
	configPaths = []string{
		".csendg.json",
		HomeDir() + "/.config/csendg/config.json",
	}
)

var appConfig AppConfig

func init() {
	// 設定ファイルが存在したら読み込み
	for _, p := range configPaths {
		if Exists(p) {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(b, &appConfig); err != nil {
				fmt.Fprintln(os.Stderr, "[err]設定ファイル("+p+")の書式が不正です。")
				panic(err)
			}
			return
		}
	}
	panic("設定ファイルは必須")
}

func main() {
	config := oauth2.Config{
		ClientID:     appConfig.ClientID,
		ClientSecret: appConfig.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",          //今回はリダイレクトしないためこれ
		Scopes:       []string{"https://mail.google.com/"}, //必要なスコープを追加
	}

	expiry, _ := time.Parse("2006-01-02", "2017-07-11")
	token := oauth2.Token{
		AccessToken:  appConfig.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: appConfig.RefreshToken,
		Expiry:       expiry,
	}

	client := config.Client(oauth2.NoContext, &token)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	// r, err := srv.Users.Labels.List("me").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to get labels. %v", err)
	// }
	//
	// if len(r.Labels) > 0 {
	// 	fmt.Print("Labels:\n")
	// 	for _, l := range r.Labels {
	// 		fmt.Printf("- %s\n", l.Name)
	// 	}
	// } else {
	// 	fmt.Print("No label found.")
	// }
	//
	temp := []byte("From: 'me'\r\n" +
		"To: jiroron666@gmail.com\r\n" +
		"Subject: TestSubject\r\n" +
		"\r\n" + "TestBody")

	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("%v", err)
	}
}
