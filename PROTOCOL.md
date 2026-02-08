# Claude Code CLI stream-json プロトコル

## シナリオ

### シンプルなテキスト応答の基本フロー

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Hello!"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","result":"Hello!"}` |

### 単一のBashツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"The command printed: tool-use-test-output"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### 2段階の連続ツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Both commands completed successfully."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### テキストとツール呼び出しが同一レスポンスに含まれるケース

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Done. The output was: combined-test"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### 複数ツールの並列呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Both commands ran: parallel-one and parallel-two"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### 同一セッション内でのマルチターン会話

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"First answer."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Second answer."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### max_tokensで応答が打ち切られた場合の挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text"}]}}` |
| → | [result/success(is_error:true)](#resultsuccess-is_errortrue) | `{"type":"result","subtype":"success","is_error":true}` |

### 拡張思考（extended thinking）ブロックを含む応答

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(thinking)](#assistantthinking) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"Let me think about this step by step..."}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"The answer is 42."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","result":"The answer is 42."}` |

### Readツールによるファイル読み取り

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"The file contains: file-content-for-read-test"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Writeツールによるファイル作成

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"File created successfully."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Editツールによるファイル編集

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"File edited successfully."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Globツールによるファイルパターンマッチ

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Found 2 text files."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Grepツールによるコンテンツ検索

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Found the pattern in searchable.txt."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### TodoWriteツールによるタスクリスト管理

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Created a todo list with 2 items."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Read → Edit → Bashの複数ステップツールチェイン

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Chain complete: read, edited, and verified."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### 思考ブロック後にツール呼び出しが続くフロー

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"After thinking and running the command: thinking-tool-test"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### スタブAPIのリクエスト記録機能

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Done."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### NotebookEditツールによるJupyterノートブック編集

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Inserted a new cell into the notebook."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### AskUserQuestionツールの非インタラクティブモードでの挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:AskUserQuestion)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"AskUserQuestion"}]}}` |
| → | [user(tool_result:error)](#usertool_result-is_errortrue) | `{"type":"user","message":{"content":[{"type":"tool_result","is_error":true}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"You chose Go. Let me proceed with Go."}]}}` |
| → | [result/success(permission_denials)](#resultsuccess-permission_denials) | `{"type":"result","subtype":"success","permission_denials":[{"tool_name":"AskUserQuestion"}]}` |

### EnterPlanModeツールによるプランモード遷移

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:EnterPlanMode)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"EnterPlanMode"}]}}` |
| → | [system/status(plan)](#systemstatus) | `{"type":"system","subtype":"status","permissionMode":"plan"}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"I have entered plan mode. Let me explore the codebase."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### WebFetchツールのURL取得挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:WebFetch)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"WebFetch"}]}}` |
| → | [user(tool_result:error)](#usertool_result-is_errortrue) | `{"type":"user","message":{"content":[{"type":"tool_result","is_error":true}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### ツール実行失敗時のエラーハンドリング

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"The file does not exist. Let me handle this error."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### APIレベルのSSEエラーイベント受信時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [result/error_during_execution](#resulterror_during_execution) | `{"type":"result","subtype":"error_during_execution"}` |

### 複数テキストブロックを含む応答の出力形式

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"First paragraph."}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Second paragraph."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","result":"Second paragraph."}` |

### TeamCreateツールによるチーム作成

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Team created."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### TeamDeleteツールのアクティブチームなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"TeamDelete"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Handled team deletion."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### SendMessageツールのチームコンテキストなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:SendMessage)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"SendMessage"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Handled send message."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### Taskツールによるチームメイトのスポーン（生成）

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(tool_use:Task)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"Task"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

### エージェントチームのライフサイクル全体（作成→削除）をマルチターンで

| direction | message | json |
|-----------|---------|------|
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Team created successfully."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |
| ← | user | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"TeamDelete"}]}}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Team deleted."}]}}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success"}` |

## メッセージ

### system/init

CLI起動時に最初に出力されるメッセージ。セッションIDやバージョンなどの初期情報を含む。
必ずセッションの最初のメッセージとして出力され、tools フィールドに利用可能なツール一覧が含まれる。

```json
{"type":"system","subtype":"init","apiKeyStatus":"valid","cwd":"/path","sessionId":"abc","version":"1.0.0","tools":["Bash","Read"]}
```

### system/status

システム状態の変更を通知するメッセージ。
EnterPlanMode ツール実行直後に出力され、PermissionMode フィールドで現在のモードが通知される。

```json
{"type":"system","subtype":"status","permissionMode":"plan"}
```

### assistant/text

モデルのテキスト応答を含むメッセージ。
各コンテンツブロックは個別の assistant メッセージとして出力され、
複数テキストブロックがある場合はブロックごとに別メッセージになる。
text が空文字列の場合、Text フィールドは省略される（部分一致用）。

```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Hello!"}]}}
```

### assistant/tool_use

モデルがツール実行を要求するメッセージ。
各コンテンツブロックは個別の assistant メッセージとして出力される。

```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}]}}
```

### assistant/thinking

モデルの内部推論過程を含むメッセージ。
thinking ブロックは text ブロックより前に、別の assistant メッセージとして順序通り出力される。

```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"Let me think..."}]}}
```

### user/text

ユーザからCLIへ送信するメッセージ。stream-json形式でstdinに送る。
content にはユーザの指示テキストが入る。

```json
{"type":"user","message":{"role":"user","content":"say hello"}}
```

### user/tool_result

CLIがツール実行結果を報告する user メッセージ。stdout に出力される。
ContentBlock の ID に対応する tool_use_id を持つ。

```json
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]}}
```

### user/tool_result, is_error=true

ツール実行失敗またはパーミッション拒否時に出力されるエラー付きツール実行結果。

```json
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","is_error":true,"content":"Error: file not found"}]}}
```

### result/success

ターンの正常完了を示すメッセージ。Read() はこのメッセージ受信時に読み取りを終了する。
Result フィールドには最後のテキストブロックの内容が入る。
result が空文字列の場合、Result フィールドは省略される（部分一致用）。

```json
{"type":"result","subtype":"success","result":"Hello!","is_error":false,"duration_ms":1234,"duration_api_ms":1000}
```

### result/success, is_error=true

max_tokens による応答打ち切り時に出力される。Subtype は "success" だが IsError が true になる。

```json
{"type":"result","subtype":"success","is_error":true,"result":"truncated response..."}
```

### result/success, permission_denials

非インタラクティブモードでパーミッション拒否されたツールがある場合に出力される。

```json
{"type":"result","subtype":"success","permission_denials":[{"tool_name":"AskUserQuestion"}]}
```

### result/error_during_execution

APIがSSEエラーイベントを返した場合のターン終了メッセージ。
assistant メッセージは出力されず、Errors 配列にエラー種別とメッセージが格納される。

```json
{"type":"result","subtype":"error_during_execution","errors":[{"type":"overloaded_error","message":"Overloaded"}]}
```

