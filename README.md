# 

##

ログ作成

devcontainerを使った、ローカルで実行。


```bash
# LOGIN_USERNAME, LOGIN_PASSWORd
LOGIN_USERNAME=yourname LOGIN_PASSWORD=password go run main.go >> .devcontainer/logs/app.log 2>> .devcontainer/logs/app_error.log
```

## 常にセレクターで選択せよ

ブラウザから選択したいタグをクリックして、Copy Selectorとすること。
こうすると唯一に決まるため。
h1やidなど本来はhtmlの1ページ

ネットによくある自分で考えたタブを選択した程度では正しく動かない。


## うまいこと動かない時。

細かいことは考えずに手でやる方針。
癖あるので、動かないときもある。

大事なのは、ちゃんとできたかどうかの確認が取れるように実装すること。

コアなloginやログイン管理などの処理と、

ビジネスロジック分けた方がいい。
どこまで分けるべきか？
