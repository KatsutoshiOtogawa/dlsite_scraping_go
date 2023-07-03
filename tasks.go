package tasks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	// "cloud.google.com/go/logging"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	// "google.golang.org/api/option"
)

// OpenLogのCallBack関数です。このフォーマット通りに動きます。
// example: 2, 3のやり方で実装することを推奨。
//  1. vpsの場合、ローカルファイルにログを書き込む。
//     fileName := "app.log"
//     file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//     if err != nil {
//     return nil, fmt.Errorf("ログファイルを開く際にエラーが発生しました: %v", err)
//  2. devcontainer実行時や、docker logsにログを吐く場合。
//
// 特にファイルをOpenしなくて良い。ログのフォーマットも指定しなくて良い。管理が減るのでこっちの方が良い。
//
//	 何もしない関数を渡す。
//		func () (*os.File, error) {
//				return nil, nil
//		}
//
// 3. gcp, awsなどのlogging
type OpenLog func() (*os.File, error)

// CloseLogのCallBack関数です。
type CloseLog func(file *os.File) error

// WaitVisible, WaitEnableのCallback関数です。
// HeadlessならWaitEnableを選択してください。
type WaitLogic func() error

// Actionの順番に処理されるが、中で呼ぶブラウザの処理が同期的な処理とは限らないので注意。

type ScrapingTaskManager struct {
	SiteSessionCookieName string
	Width                 int64
	Height                int64
	ScreenShotLogPath     string
	ScreenShotLogPrefix   string
	SiteTopUrl            string
	LogInUrl              string
	LogOutUrl             string
	LoginPassword         string
	LoginPasswordSel      string
	LoginUsername         string
	LoginUsernameSel      string
	LoginButtonSel        string
	AgePermissionUrl      string        // 年齢認証が求められるurl
	AgePermissionSel      string        // 年齢認証が求められたときにYesを押すボタンのタグ
	AgePermissionNextSel  string        // 年齢認証が求められたときにYesを押した後に移動するページにあるSelector
	DefaultTimeSpan       time.Duration // 実行時に待つ時間のデフォルト値
	OpenLog               OpenLog
	CloseLog              CloseLog
}

// logが書けることの確認。
func IsEnableLog() (bool, error) {

	return true, nil
}

// ログをOpenする。
// gcpのcloud loggingを使う場合は
// defer CloseLog()でちゃんとファイルが閉じることを保証すること。
// この関数を渡せるように
// func OpenLog() (*os.File, error) {

// 	isLocal := false // ローカル環境かどうかの判定方法に合わせて適切な値を設定する

// 	var file *os.File

// 	if isLocal {
// 		// ローカル環境の場合、ファイルにログを書き込む
// 		fileName := "app.log"
// 		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// 		if err != nil {
// 			return nil, fmt.Errorf("ログファイルを開く際にエラーが発生しました: %v", err)
// 		}

// 		log.SetOutput(file)
// 	} else {
// 		//
// 		// GCP上の場合、Cloud Loggingにログを送信する
// 		// credential使うのはあまり良くないので、iam+serviceaccountで設定する。
// 		projectID := "your-project-id"
// 		credentialsFile := "/path/to/credentials.json"

// 		ctx := context.Background()
// 		client, err := logging.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
// 		if err != nil {
// 			return nil, fmt.Errorf("Cloud Loggingクライアントの作成に失敗しました: %v", err)
// 		}

// 		logger := client.Logger("my-log")
// 		log.SetOutput(logger.Writer(logging.Info))
// 	}

// 	return file, nil

// }

// func CloseLog(file *os.File) error {

// 	isLocal := false // ローカル環境かどうかの判定方法に合わせて適切な値を設定する

// 	if isLocal {
// 		err := file.Close()
// 		if err != nil {
// 			return fmt.Errorf("Cloud Loggingクライアントの作成に失敗しました: %v", err)
// 		}
// 	}

// 	// gcpなどでは常にnil

// 	return nil
// }

// 年齢認証通過後か判定
func (s ScrapingTaskManager) IsAgeVerificationTasks(targetCookieName string, targetCookieValue string, valid *bool) chromedp.Tasks {

	file, _ := s.OpenLog()

	defer s.CloseLog(file)

	// 判定が終わるまでunkownなのでnilを入れる
	valid = nil
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				log.Println("cookieが取得できませんでした。", err)
				return err
			}

			for _, cookie := range cookies {
				if cookie.Name == targetCookieName && cookie.Value == targetCookieValue {
					*valid = true
					break
				}
			}

			if valid == nil {
				*valid = false
			}

			return nil
		}),
	}
}

// 常にエラーになることが保証されているタスク
// エラーに入れたい文字列を代入してください。
func (s ScrapingTaskManager) ErrorTask(v string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			return errors.New(v)
		}),
	}
}

// セッションが有効かどうかの確認を行う。
func (s ScrapingTaskManager) IsSessionVerificationTasks(valid *bool) chromedp.Tasks {

	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	// 判定が終わるまでunkownなのでnilを入れる
	valid = nil
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				log.Println("cookieが取得できませんでした。", err)
				return err
			}

			for _, cookie := range cookies {
				if cookie.Name == s.SiteSessionCookieName {
					*valid = true
					break
				}
			}

			if valid == nil {
				*valid = false
			}

			return nil
		}),
	}
}

// ウィンドウのサイズを調整する
func (s ScrapingTaskManager) EmulateViewportTasks(width int64, height int64) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// ウィンドウサイズを指定（オプション）
			err := chromedp.Run(ctx, chromedp.EmulateViewport(width, height))
			if err != nil {
				log.Println("ウィンドウサイズの変更ができませんでした。", err)
				return err
			}

			return nil
		}),
	}
}

// url全体を取得する
func (s ScrapingTaskManager) LocationHrefTasks(href *string) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {

			err := chromedp.Run(ctx,
				chromedp.EvaluateAsDevTools("window.location.href", href),
			)
			if err != nil {
				log.Println("hrefが取得できませんでした。", err)
				return err
			}

			return nil
		}),
	}
}

// ウィンドウのサイズを取得する
func (s ScrapingTaskManager) ViewSizeTasks(width *int64, height *int64) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// ウィンドウサイズを指定（オプション）
			err := chromedp.Run(ctx,
				chromedp.EvaluateAsDevTools("window.innerHeight", height),
				chromedp.EvaluateAsDevTools("window.innerWidth", width),
			)
			if err != nil {
				log.Println("ウィンドウサイズが取得できませんでした。。", err)
				return err
			}

			return nil
		}),
	}
}

// logのためにscreenShotをとる。
// args:
//
//			sel h1, div1
//		 fileNameFormat prefixの後のファイル名。自由形式。空の文字列でも良い。
//	  fileExtension 現時点ではpng固定。
func (s ScrapingTaskManager) TakeScreenShotLogTasks(sel interface{}, logFileName string, fileExtension string) chromedp.Tasks {
	// スクリーンショットの名称指定。
	currentTime := time.Now().Format(s.ScreenShotLogPrefix)

	fileName := filepath.Join(s.ScreenShotLogPath, fmt.Sprintf("%s%s.%s", currentTime, logFileName, fileExtension))
	return s.TakeScreenShotTasks(sel, fileName)
}

// ウィンドウのサイズを調整する
// screenShotをとる。
// args:
//
//		sel h1, div1
//	 fileName パスと拡張子まで含めた
func (s ScrapingTaskManager) TakeScreenShotTasks(sel interface{}, fileName string) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			var imageBuf []byte
			// スクリーンショットを取得
			// スクリーンショットの名称指定。
			err := chromedp.Screenshot(sel, &imageBuf, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			if err != nil {
				log.Println("スクリーンショットが取得できませんでした。", err)
				return err
			}

			// スクリーンショットをファイルに保存
			err = os.WriteFile(fileName, imageBuf, 0640)
			if err != nil {
				log.Println("スクリーンショットをファイルに保存できませんでした。", err)
				return err
			}

			return nil
		}),
	}
}

// キー入力を行う。
func (s ScrapingTaskManager) SendKeysTasks(sel interface{}, v string, t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.SendKeys(sel, v).Do(ctx)
			if err != nil {
				log.Println("キー入力ができませんでした。", err)
				// logなど。
			}

			return nil
		}),
		s.WaitTasks(waitTime),
	}
}

// サイトのトップページに移動する
func (s ScrapingTaskManager) MoveTopPageTasks() chromedp.Tasks {
	return s.MovePageTasks(s.SiteTopUrl)
}

// サイトの特定のページに移動する
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) MovePageTasks(url string, t ...time.Duration) chromedp.Tasks {
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}

	return chromedp.Tasks{
		chromedp.Navigate(url),
		s.WaitTasks(waitTime),
	}
}

// サイトにログインする
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) LoginSiteTasks(t ...time.Duration) chromedp.Tasks {

	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}

	return chromedp.Tasks{
		// ログイン状態かどうかの検証。
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	var valid bool
		// 	err := s.IsSessionVerificationTasks(&valid).Do(ctx)
		// 	if err != nil {
		// 		log.Println("ログイン済みかどうか確認できませんでした。", err)
		// 		return err
		// 	}

		// 	if valid {
		// 		return errors.New("ログイン状態で、ログインを呼び出そうとしました。")
		// 	}
		// 	return nil
		// }),
		s.MovePageTasks(s.LogInUrl),
		s.TakeScreenShotLogTasks("html", "login", "png"),
		s.SendKeysTasks(s.LoginUsernameSel, s.LoginUsername, waitTime),
		s.SendKeysTasks(s.LoginPasswordSel, s.LoginPassword, waitTime),
		s.ClickTasks(s.LoginButtonSel, waitTime),
		s.WaitTasks(waitTime),
	}
}

// サイトにログアウトする
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) LogoutTasks(t ...time.Duration) chromedp.Tasks {
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}

	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			var valid bool
			err := s.IsSessionVerificationTasks(&valid).Do(ctx)
			if err != nil {
				log.Println("ログイン済みかどうか確認できませんでした。", err)
				return err
			}

			if !valid {
				return errors.New("未ログイン状態で、ログアウトを呼び出そうとしました。")
			}
			return nil
		}),
		// ボタンをクリックするかどうかはサイトによる。
		s.MovePageTasks(s.LogOutUrl),
		// s.ClickTasks("#passwordNext", waitTime),
		s.WaitTasks(waitTime),
	}
}

// 要素をクリックする
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) ClickTasks(sel interface{}, t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.Click(sel).Do(ctx)
			if err != nil {
				log.Println("クリックできませんでした。", err)
				return err
			}
			return nil
		}),
		s.WaitTasks(waitTime),
	}
}

// 処理を待つのに使う
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) WaitTasks(t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.Sleep(waitTime).Do(ctx)
			if err != nil {
				log.Println("待てませんでした。", err)
				return err
			}
			log.Printf("%v待ちました。", waitTime)
			return nil
		}),
	}
}

// テキストを取得するのに使う。
// Selectorに合致するものが複数あると正しく動かないので注意。
// 主に正しく実行されているかなどの検査に使う。
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) TextContentTasks(sel interface{}, v *string, t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{

		// s.WaitVisibleTasks(sel),
		s.WaitEnableTasks(sel),

		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.TextContent(sel, v).Do(ctx)

			if err != nil {
				log.Println("textContentを取得できませんでした。", err)
				return err
			}
			log.Printf("textContentの値は%sです", *v)
			return nil
		}),
		// chromedp.TextContent(sel, v),
		s.WaitTasks(waitTime),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	err := chromedp.Sleep(waitTime).Do(ctx)
		// 	if err != nil {
		// 		log.Println("待てませんでした。", err)
		// 		return err
		// 	}
		// 	log.Printf("%v待ちました。", waitTime)
		// 	return nil
		// }),
	}
}

// 要素が見えるのを待つ。Headlessなら永遠に表示されないので、使わない。
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) WaitVisibleTasks(sel interface{}, t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	// 何を待っているかのログを数秒置きに出す。
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.WaitVisible(sel).Do(ctx)
			if err != nil {
				log.Println("要素が見えるのを待てませんでした。", err)
				return err
			}
			return nil
		}),
		s.WaitTasks(waitTime),
	}
}

// 要素が使えるようになるのを待つ。
// t 待つのに使う。指定しない場合は、デフォルトの時間が使われる。配列の最初にあるものしか使われない。
func (s ScrapingTaskManager) WaitEnableTasks(sel interface{}, t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.WaitEnabled(sel).Do(ctx)
			if err != nil {
				log.Println("要素が使えるようになるのを待てませんでした。", err)
				return err
			}
			return nil
		}),
		s.WaitTasks(waitTime),
	}
}

// 年齢認証を突破する
func (s ScrapingTaskManager) AgeVerificationTasks(t ...time.Duration) chromedp.Tasks {
	file, _ := s.OpenLog()

	defer s.CloseLog(file)
	var waitTime time.Duration
	if len(t) == 0 {
		waitTime = s.DefaultTimeSpan
	} else {
		waitTime = t[0]
	}
	return chromedp.Tasks{
		// 年齢認証が必要な場所に移動する。
		s.MovePageTasks(s.AgePermissionUrl),
		// chromedp.Navigate(url),
		// chromedp.Sleep(2 * time.Second), // ページの読み込みを待つために適切な時間を設定してください

		// chromedp.ActionFunc(func(ctx context.Context) error {

		// 	// 年齢確認ボタンをクリック
		// 	err := chromedp.Click("#age_check_button").Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	return nil
		// }),
		// 年齢確認ボタンをクリック
		s.ClickTasks(s.AgePermissionSel, waitTime),
		// chromedp.ActionFunc(func(ctx context.Context) error {

		// 	// 成人向けコンテンツのページが読み込まれるのを待つ
		// 	err := chromedp.WaitVisible("#adult_contents_div").Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	return nil
		// }),

		// 成人向けコンテンツページが読み込まれるのを待つ。
		// たいてい自動的に成人向けに移動されるため。
		s.WaitVisibleTasks(s.AgePermissionNextSel, waitTime),

		// chromedp.ActionFunc(func(ctx context.Context) error {

		// 	// 年齢確認ボタンをクリック
		// 	err := chromedp.Click("#age_check_button").Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	// 成人向けコンテンツのページが読み込まれるのを待つ
		// 	err = chromedp.WaitVisible("#adult_contents_div").Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	return nil
		// }),
		// takeScreenShot(`h1`, screenshotPath,)
	}
}
