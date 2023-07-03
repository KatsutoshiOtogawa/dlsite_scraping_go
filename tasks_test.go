package tasks

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// ページの遷移機能が正しく動作するか確認。
func TestMoveTopPageTasks(t *testing.T) {

	ctx, cancel := chromedp.NewContext(context.Background())

	// ctx, cancel := chromedp.NewExecAllocator(context.Background(),
	// 	chromedp.Flag("headless", true), // ヘッドレスモードを有効にする
	// )

	defer cancel()

	// タイムアウトの設定
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	//　これが小さすぎるとcontext deadline exceedになる。
	ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	num, err := strconv.Atoi(os.Getenv("DEFAULT_TIME_SPAN"))
	if err != nil {
		t.Errorf("整数に変換できませんでした。 MoveTopPageTasks() = %v", err)
		return
	}
	taskManager := ScrapingTaskManager{
		SiteSessionCookieName: os.Getenv("SITE_SESSION_COOKIE"),
		SiteTopUrl:            os.Getenv("SITE_TOP_URL"),
		DefaultTimeSpan:       time.Duration(num) * time.Second,
		ScreenShotLogPath:     os.Getenv("SCREENSHOT_LOG_PATH"),
		ScreenShotLogPrefix:   os.Getenv("SCREENSHOT_LOG_PREFIX"),
		LogInUrl:              os.Getenv("LOGIN_URL"),
		LoginUsername:         os.Getenv("LOGIN_USERNAME"),
		LoginUsernameSel:      os.Getenv("LOGIN_USERNAME_SEL"),
		LoginPasswordSel:      os.Getenv("LOGIN_PASSWORD_SEL"),
		LoginButtonSel:        os.Getenv("LOGIN_BUTTON_SEL"),
		LogOutUrl:             os.Getenv("LOGOUT_URL"),

		OpenLog: func() (*os.File, error) {

			return nil, nil
		},
		CloseLog: func(file *os.File) error {
			return nil
		},
	}

	type args struct {
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "MoveTopPage",
			args: args{},
			want: "Select Language",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var text string
			err := chromedp.Run(ctx,
				taskManager.MoveTopPageTasks(),
				// 移動した後のテキストを取得する
				taskManager.TextContentTasks("#locale_setting_title", &text),
				taskManager.TakeScreenShotLogTasks("#locale_setting_title", "", "png"),
			)
			// エラーが出た。
			if err != nil {
				log.Fatal(err)
				t.Errorf("MoveTopPageTasks() = %v", err)
			}

			if text != tt.want {
				t.Errorf("MoveTopPageTasks() = %s, want %s", text, tt.want)
			}
		})
	}

}

// サイトログイン,ログアウト機能が正しく動作するか確認。
func TestLoginLogoutTasks(t *testing.T) {

	ctx, cancel := chromedp.NewContext(context.Background())

	// ctx, cancel := chromedp.NewExecAllocator(context.Background(),
	// 	chromedp.Flag("headless", true), // ヘッドレスモードを有効にする
	// )

	defer cancel()

	// タイムアウトの設定
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	//　これが小さすぎるとcontext deadline exceedになる。
	ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	num, err := strconv.Atoi(os.Getenv("DEFAULT_TIME_SPAN"))
	if err != nil {
		t.Errorf("整数に変換できませんでした。 MoveTopPageTasks() = %v", err)
		return
	}
	taskManager := ScrapingTaskManager{
		SiteSessionCookieName: os.Getenv("SITE_SESSION_COOKIE"),
		SiteTopUrl:            os.Getenv("SITE_TOP_URL"),
		DefaultTimeSpan:       time.Duration(num) * time.Second,
		ScreenShotLogPath:     os.Getenv("SCREENSHOT_LOG_PATH"),
		ScreenShotLogPrefix:   os.Getenv("SCREENSHOT_LOG_PREFIX"),
		LogInUrl:              os.Getenv("LOGIN_URL"),
		LoginUsername:         os.Getenv("LOGIN_USERNAME"),
		LoginUsernameSel:      os.Getenv("LOGIN_USERNAME_SEL"),
		LoginPasswordSel:      os.Getenv("LOGIN_PASSWORD_SEL"),
		LoginButtonSel:        os.Getenv("LOGIN_BUTTON_SEL"),
		LogOutUrl:             os.Getenv("LOGOUT_URL"),

		OpenLog: func() (*os.File, error) {

			return nil, nil
		},
		CloseLog: func(file *os.File) error {
			return nil
		},
	}

	type args struct {
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "MoveTopPage",
			args: args{},
			// login状態が無効か、有効かを表す
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var valid bool
			err := chromedp.Run(ctx,
				taskManager.LoginSiteTasks(),
				// Loginした後にsessionが有効になったかチェックする。
				taskManager.IsSessionVerificationTasks(&valid),
			)
			// エラーが出た。
			if err != nil {
				log.Fatal(err)
				t.Errorf("TestLoginTasks() = %v", err)
			}

			if valid == tt.want {
				t.Errorf("TestLoginTasks() = %v, want %v", valid, tt.want)
			}

			err = chromedp.Run(ctx,
				taskManager.LogoutTasks(),
				// Logoutした後にsessionが無効になったかチェックする。
				taskManager.IsSessionVerificationTasks(&valid),
			)
			// エラーが出た。
			if err != nil {
				log.Fatal(err)
				t.Errorf("TestLoginTasks() = %v", err)
			}

			if valid != tt.want {
				t.Errorf("TestLoginTasks() = %v, want %v", valid, tt.want)
			}
		})
	}

}
