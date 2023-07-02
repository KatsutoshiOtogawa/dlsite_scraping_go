package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

// 移動したいコンテンツへのurl作成ロジック
func makeMainContentsUrl(siteUrl string) string {

	contents := os.Getenv("CONTENTS")
	// DLsiteのURL（成人向けコンテンツに直接リンクするページのURLを指定してください）
	contentsUrl := filepath.Join(siteUrl, contents)

	return contentsUrl
}

func main() {
	// コンテキストの作成
	ctx, cancel := chromedp.NewContext(context.Background())

	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.Flag("headless", true), // ヘッドレスモードを有効にする
	)

	defer cancel()

	// タイムアウトの設定
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	siteUrl := os.Getenv("SITE_TOP_URL")

	// Chromeを起動してDLsiteにアクセス
	err := chromedp.Run(ctx, MoveTopPageTasks(siteUrl))
	if err != nil {
		log.Fatal(err)
	}

	// 年齢認証確認後か確認
	var ageVerification bool

	targetCookieName := os.Getenv("ADULT_SESSION_KEY")
	targetCookieValue := os.Getenv("ADULT_SESSION_VALUE")
	err = chromedp.Run(ctx, IsAgeVerificationTasks(targetCookieName, targetCookieValue, &ageVerification))
	if err != nil {
		log.Fatal(err)
	}

	// スクリーンショット保存先のパス
	screenshotPath := "screenshot.png"

	contentsUrl := makeMainContentsUrl(siteUrl)

	// currentTime := time.Now().Format("2006-01-02 15:04:05.000000")

	// format例を入れるといい感じに表示してくれる。
	fileNameFormat := "2006-01-02_15:04:05.000000"
	// Chromeを起動して目的のコンテンツへDLsiteにアクセス
	err = chromedp.Run(ctx, AgeVerificationTasks(contentsUrl, screenshotPath, fileNameFormat))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("スクリーンショットを保存しました:", screenshotPath)

	// 目的のCookieが存在するかどうかを確認
	sessionVerification := false
	targetCookieName := os.Getenv("SITE_SESSION_COOKIE")
	// "session_state"
	err = chromedp.Run(ctx, IsSessionVerificationTasks(targetCookieName, &sessionVerification))
	if err != nil {
		log.Fatal(err)
	}

}

func MoveTopPageTasks(siteUrl string) {
	panic("unimplemented")
}
