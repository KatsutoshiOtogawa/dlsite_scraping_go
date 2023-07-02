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
		SessionCookieName:   os.Getenv("SITE_SESSION_COOKIE"),
		SiteTopUrl:          os.Getenv("SITE_TOP_URL"),
		DefaultTimeSpan:     time.Duration(num) * time.Second,
		ScreenShotLogPath:   os.Getenv("SCREENSHOT_LOG_PATH"),
		ScreenShotLogPrefix: os.Getenv("SCREENSHOT_LOG_PREFIX"),

		OpenLog: func() (*os.File, error) {

			return nil, nil
		},
		CloseLog: func(file *os.File) error {
			return nil
		},
	}

	// err := chromedp.Run(ctx, taskManager.MoveTopPageTasks())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // 移動した後のテキストを取得する
	// var text string
	// err = chromedp.Run(ctx, taskManager.TextContentTasks("h1", text))
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
			// if got := Sub(tt.args.i, tt.args.j); got != tt.want {
			// 	t.Errorf("Sub() = %v, want %v", got, tt.want)
			// }
			// err := chromedp.Run(ctx, taskManager.MoveTopPageTasks())
			// if err != nil {
			// 	t.Errorof(err)
			// 	log.Fatal(err)
			// }
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

			// if got := Sub(tt.args.i, tt.args.j); got != tt.want {
			// 	t.Errorf("Sub() = %v, want %v", got, tt.want)
			// }
		})
	}

}
