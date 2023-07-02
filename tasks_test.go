package tasks_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// テストケースを定義します
var testCases = []struct {
	input    int
	expected int
}{
	{2, 4},   // 2を入力した場合、4が期待値
	{0, 0},   // 0を入力した場合、0が期待値
	{-1, -2}, // -1を入力した場合、-2が期待値
}

// テスト対象の関数を定義します
func Double(x int) int {
	return x * 2
}

// テスト関数を定義します
func TestDouble(t *testing.T) {
	// 各テストケースを順に実行します
	for _, tc := range testCases {
		// テストケースを実行し、結果を取得します
		result := Double(tc.input)

		// 結果と期待値を比較し、一致しない場合にエラーメッセージを表示します
		if result != tc.expected {
			t.Errorf("Double(%d) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

func TestMovePageTasks(t *testing.T) {

	ctx, cancel := chromedp.NewContext(context.Background())

	// ctx, cancel := chromedp.NewExecAllocator(context.Background(),
	// 	chromedp.Flag("headless", true), // ヘッドレスモードを有効にする
	// )

	defer cancel()

	// タイムアウトの設定
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var taskManager ScrapingTaskManager

	taskManager.MoveTopPageTasks(siteUrl)

	err := chromedp.Run(ctx, MoveTopPageTasks(siteUrl))
	if err != nil {
		log.Fatal(err)
	}
}
