# 

##

ログ作成

devcontainerを使った、ローカルで実行。


```bash
go run main.go >> .devcontainer/logs/app.log 2>> .devcontainer/logs/app_error.log
```


## うまいこと動かない時。

細かいことは考えずに手でやる方針。
癖あるので、動かないときもある。

大事なのは、ちゃんとできたかどうかの確認が取れるように実装すること。
