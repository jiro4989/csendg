package main

import (
	"bufio"
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
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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

	args := os.Args
	toAddr := args[1]
	filePath := args[2]

	mf, err := ReadMailFile(filePath)
	if err != nil {
		panic(err)
	}

	// Subjectの文字化けのため
	msgStr := "From: 'me'\r\n" +
		"To: " + toAddr + "\r\n" +
		"Subject: " + mf.Title + "\r\n" +
		"\r\n" + mf.Body

	iso2022jpMsg, err := toISO2022JP(msgStr)
	if err != nil {
		panic(err)
	}
	msg := []byte(iso2022jpMsg)

	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(msg)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("%v", err)
	}
}

type MailFile struct {
	FilePath, Title, Body string
}

func ReadMailFile(p string) (MailFile, error) {
	f, err := os.Open(p)
	if err != nil {
		return MailFile{}, err
	}
	defer f.Close()

	var (
		ln    int
		lines []string
		mf    = MailFile{FilePath: p}
		sc    = bufio.NewScanner(f)
	)

	for sc.Scan() {
		line := sc.Text()
		switch ln {
		case 0:
			mf.Title = line
		case 1:
			ln++
			continue
		default:
			lines = append(lines, line)
		}
		ln++
	}
	if err := sc.Err(); err != nil {
		return MailFile{}, err
	}
	mf.Body = strings.Join(lines, "\n")
	return mf, nil
}

// Convert UTF-8 to ISO2022JP
func toISO2022JP(str string) ([]byte, error) {
	reader := strings.NewReader(str)
	transformer := japanese.ISO2022JP.NewEncoder()
	return ioutil.ReadAll(transform.NewReader(reader, transformer))
}
