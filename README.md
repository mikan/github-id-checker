github-id-checker
=================

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/mikan/github-id-checker)

GitHub の ID が Organization の運用ルールに適合しているかチェックし、適合している場合は登録申請を行えるようにするアプリケーションです。

## 必要な環境変数

予め以下の環境変数を刺しておく必要があります:

| Var | 説明 |
| --- | ---- |
| ORG | 対象 GitHub Organization |
| KEYWORD | メールアドレスに含まれるべき文字列 |
| POLICY_URL | プライバシーポリシーのリンク先 URL |
| CLIENT_ID | GitHub OAuth の Client ID |
| CLIENT_SECRET | GitHub OAuth の Client Secret |
| WEBHOOK | 通知先 Webhook URL |

GitHub の Client ID/Secret を得るには、 Organization の [Settings] から [OAuth Apps] を選び [New OAuth App] ボタンから必要事項を記入して登録することで得られます。
登録の際の [Authorization callback URL] は、デプロイした URL の後ろに `/github` を加えた URL を指定します。

名前の入力チェックは以下の条件を全て満たすものとしています:

- Trim 処理後の文字数が 3 文字以上
- Trim 処理後に半角スペースを含む

メールの通知を行う場合は、以下の設定が必要です:

| Var | 説明 |
| --- | ---- |
| SEND_FROM | エントリーの通知元メールアドレス |
| SEND_TO | エントリーの通知先メールアドレス |
| SMTP_USER | SMTP ユーザー |
| SMTP_PASSWORD | SMTP パスワード |
| SMTP_SERVER | SMTP サーバー |
| SMTP_PORT | SMTP ポート |


## 開発

### 概要

実装言語は Go で、標準ライブラリと公式ライブラリの `golang.org/x/oauth2` のみを用いて開発しています。

UI 部分は素の HTML + CSS です。

### コードスタイル

Go コードは Go 標準の `gofmt` でフォーマットしてください。
`import` 文の編成機能のある上位コマンド `goimports` の利用を推奨します。

それ以外のコード (shell, html, css) は GoLand (IntelliJ) の標準フォーマッターを推奨します。

## Author

- [mikan](https://github.com/mikan)
