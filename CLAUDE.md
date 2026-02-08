Claude Code CLIの `--output-format stream-json` プロトコルを観測ベースで解析し、回帰テストとして記録するプロジェクト。

# 目的

- 解析対象は Claude Code CLI（`claude` コマンド）の stream-json I/O のみ。
- テストコード → Claude Code CLI → Stub API という構成で、CLIの入出力を観測する。
  - input: `--input-format stream-json` でstdinへ送るJSONL
  - output: `--output-format stream-json` でstdoutから受け取るJSONL
- 観測結果をGoテストとして記録し、CLIの破壊的変更を検知する回帰テストとする。

# 構成

```
protocol_test.go     -- 観測・回帰テスト本体
utils/
  stub_api.go        -- Anthropic Messages API のSSEストリーミングスタブ
  harness.go         -- CLIプロセス管理とアサーションユーティリティ
```

# テストの書き方

## 基本パターン

1. `utils.StubAPIServer` に返すSSEレスポンス列を設定し、`Start()` する
2. `utils.NewSession(t, stub.URL())` でCLIプロセスを起動する
3. `s.Send(...)` でstream-json形式のユーザメッセージをstdinに送る
4. `s.Read()` でstdoutから `result` メッセージまでを読み取る
5. `utils.AssertOutput(t, output, patterns...)` で部分一致アサーションする

## レスポンスヘルパー

- `utils.TextResponse(text)` — テキスト応答のSSEイベント列を生成する
- `utils.ToolUseResponse(toolID, toolName, input)` — ツール呼び出しのSSEイベント列を生成する（stop_reason: "tool_use"）

ツール利用シナリオでは `Responses` に複数のレスポンスを設定する。リクエスト順にレスポンスが消費され、超過分は最後のレスポンスが繰り返される。

## AssertOutput

パターンはJSON部分一致。指定したフィールドが実際のメッセージに含まれていればマッチする。出力中でパターンを順序通りに探し、マッチしないメッセージはスキップされる。

# 前提・制約

- CLIへの入力は必ずstdin経由、出力は必ずstdout経由で行う。
- CLIは `--dangerously-skip-permissions --verbose --no-session-persistence` 付きで起動される。
- `ANTHROPIC_BASE_URL` 環境変数でスタブAPIに向ける。
- 観測した形式を推測で固定しない。未知フィールドは無視してアサーションは部分一致で行う。
- 利用言語はGo。`protocol_test.go` がメインの成果物。
