# Claude Code CLI stream-json プロトコル

## シナリオ

### シンプルなテキスト応答の基本フロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success","result":"Hello!"}` |

### 単一のBashツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The command printed: tool-use-test-output"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### 2段階の連続ツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands completed successfully."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### テキストとツール呼び出しが同一レスポンスに含まれるケース

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Done. The output was: combined-test"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### 複数ツールの並列呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands ran: parallel-one and parallel-two"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### 同一セッション内でのマルチターン会話

| direction | message | json |
|-----------|---------|------|
| ← 1 | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First answer."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |
| ← 2 | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second answer."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### max_tokensで応答が打ち切られた場合の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text"}]}}` |
| → | [result/success(is_error:true)](#message) | `{"type":"result","subtype":"success","is_error":true}` |

### 拡張思考（extended thinking）ブロックを含む応答

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(thinking)](#message) | `{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think about this step by step..."}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The answer is 42."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success","result":"The answer is 42."}` |

### Readツールによるファイル読み取り

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The file contains: file-content-for-read-test"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Writeツールによるファイル作成

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"File created successfully."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Editツールによるファイル編集

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"File edited successfully."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Globツールによるファイルパターンマッチ

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Found 2 text files."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Grepツールによるコンテンツ検索

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Found the pattern in searchable.txt."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### TodoWriteツールによるタスクリスト管理

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Created a todo list with 2 items."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Read → Edit → Bashの複数ステップツールチェイン

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Chain complete: read, edited, and verified."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### 思考ブロック後にツール呼び出しが続くフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"After thinking and running the command: thinking-tool-test"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### スタブAPIのリクエスト記録機能

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Done."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### NotebookEditツールによるJupyterノートブック編集

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Inserted a new cell into the notebook."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### AskUserQuestionツールの非インタラクティブモードでの挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:AskUserQuestion)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"AskUserQuestion"}]}}` |
| → | [user(tool_result:error)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result","is_error":true}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"You chose Go. Let me proceed with Go."}]}}` |
| → | [result/success(permission_denials)](#message) | `{"type":"result","subtype":"success","permission_denials":[{"tool_name":"AskUserQuestion"}]}` |

### EnterPlanModeツールによるプランモード遷移

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:EnterPlanMode)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"EnterPlanMode"}]}}` |
| → | [system/status(plan)](#message) | `{"type":"system","subtype":"status","permissionMode":"plan"}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"I have entered plan mode. Let me explore the codebase."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### WebFetchツールのURL取得挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:WebFetch)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"WebFetch"}]}}` |
| → | [user(tool_result:error)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result","is_error":true}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### ツール実行失敗時のエラーハンドリング

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The file does not exist. Let me handle this error."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### APIレベルのSSEエラーイベント受信時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [result/error_during_execution](#message) | `{"type":"result","subtype":"error_during_execution"}` |

### 複数テキストブロックを含む応答の出力形式

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First paragraph."}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second paragraph."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success","result":"Second paragraph."}` |

### TeamCreateツールによるチーム作成

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### TeamDeleteツールのアクティブチームなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamDelete)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"TeamDelete"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled team deletion."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### SendMessageツールのチームコンテキストなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:SendMessage)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled send message."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### Taskツールによるチームメイトのスポーン（生成）

| direction | message | json |
|-----------|---------|------|
| ← | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(tool_use:Task)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Task"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

### エージェントチームのライフサイクル全体（作成→削除）をマルチターンで

| direction | message | json |
|-----------|---------|------|
| ← 1 | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [system/init](#message) | `{"type":"system","subtype":"init"}` |
| → | [assistant(tool_use:TeamCreate)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created successfully."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |
| ← 2 | [user](#message) | `{"type":"user","message":{"role":"user","content":"..."}}` |
| → | [assistant(tool_use:TeamDelete)](#message) | `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"TeamDelete"}]}}` |
| → | [user(tool_result)](#message) | `{"type":"user","message":{"content":[{"type":"tool_result"}]}}` |
| → | [assistant(text)](#message) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team deleted."}]}}` |
| → | [result/success](#message) | `{"type":"result","subtype":"success"}` |

## スキーマ

### enum 定数

#### MessageType

MessageType はプロトコルメッセージの type フィールドの値。

| 定数名 | 値 | 説明 |
|--------|-----|------|
| `TypeSystem` | `"system"` | TypeSystem はCLIのライフサイクルイベント（初期化、状態変更）。 |
| `TypeAssistant` | `"assistant"` | TypeAssistant はモデルからの応答メッセージ。 |
| `TypeUser` | `"user"` | TypeUser はユーザ入力またはツール実行結果。 |
| `TypeResult` | `"result"` | TypeResult はターンの完了通知。 |

#### Subtype

Subtype はプロトコルメッセージの subtype フィールドの値。
system メッセージと result メッセージで使用される。

| 定数名 | 値 | 説明 |
|--------|-----|------|
| `SubtypeInit` | `"init"` | SubtypeInit はCLI起動時の初期化メッセージ (type=system)。 |
| `SubtypeStatus` | `"status"` | SubtypeStatus はシステム状態の変更通知 (type=system)。 |
| `SubtypeSuccess` | `"success"` | SubtypeSuccess はターンの正常完了 (type=result)。 |
| `SubtypeErrorDuringExecution` | `"error_during_execution"` | SubtypeErrorDuringExecution はAPIエラーによるターン終了 (type=result)。 |

#### ContentBlockType

ContentBlockType はコンテンツブロックの type フィールドの値。
assistant メッセージと user メッセージの content 配列内で使用される。

| 定数名 | 値 | 説明 |
|--------|-----|------|
| `BlockText` | `"text"` | BlockText はテキスト応答ブロック。モデルの自然言語応答を含む。 |
| `BlockToolUse` | `"tool_use"` | BlockToolUse はツール呼び出しブロック。モデルがツール実行を要求する。 |
| `BlockThinking` | `"thinking"` | BlockThinking は拡張思考ブロック。モデルの内部推論過程を含む。 |
| `BlockToolResult` | `"tool_result"` | BlockToolResult はツール実行結果ブロック。CLIがツール実行結果を報告する。 |

### 型定義

#### Message

Message はプロトコルメッセージの統合表現。
スキーマドキュメントの定義ソースであり、テストのアサーションパターンにも使用する。
全フィールドに omitempty を付与しており、[MustJSON] でJSON化すると
設定されたフィールドのみが出力される（部分一致パターンとして機能する）。

アサーションパターンの構築には NewMessageXxx コンストラクタ関数を使用する。
各コンストラクタにメッセージ種別ごとのドキュメントが記載されている。

| フィールド | JSON キー | 説明 |
|-----------|-----------|------|
| `Type` | `type` | メッセージ種別 |
| `Subtype` | `subtype` | サブタイプ (system, result で使用) |
| `Message` | `message` | メッセージ本体 (assistant, user で使用) |
| `Result` | `result` | 最後のテキストブロックの内容 (result で使用) |
| `IsError` | `is_error` | エラー時 true (result で使用) |
| `PermissionMode` | `permissionMode` | パーミッションモード (system/status で使用) |
| `PermissionDenials` | `permission_denials` | パーミッション拒否されたツール (result で使用) |
| `Errors` | `errors` | APIエラー配列 (result で使用) |

##### system/init

CLI起動時に最初に出力されるメッセージ。セッションIDやバージョンなどの初期情報を含む。
必ずセッションの最初のメッセージとして出力され、tools フィールドに利用可能なツール一覧が含まれる。
Example:
```json
{"type":"system","subtype":"init","apiKeyStatus":"valid","cwd":"/path","sessionId":"abc","version":"1.0.0","tools":["Bash","Read"]}
```

##### system/status

システム状態の変更を通知するメッセージ。
EnterPlanMode ツール実行直後に出力され、PermissionMode フィールドで現在のモードが通知される。
Example:
```json
{"type":"system","subtype":"status","permissionMode":"plan"}
```

##### assistant/text

モデルのテキスト応答を含むメッセージ。
各コンテンツブロックは個別の assistant メッセージとして出力され、
複数テキストブロックがある場合はブロックごとに別メッセージになる。
text が空文字列の場合、Text フィールドは省略される（部分一致用）。
Example:
```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Hello!"}]}}
```

##### assistant/tool_use

モデルがツール実行を要求するメッセージ。
各コンテンツブロックは個別の assistant メッセージとして出力される。
Example:
```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}]}}
```

##### assistant/thinking

モデルの内部推論過程を含むメッセージ。
thinking ブロックは text ブロックより前に、別の assistant メッセージとして順序通り出力される。
Example:
```json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"Let me think..."}]}}
```

##### user/tool_result

CLIがツール実行結果を報告する user メッセージ。stdout に出力される。
ContentBlock の ID に対応する tool_use_id を持つ。
Example:
```json
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]}}
```

##### user/tool_result, is_error=true

ツール実行失敗またはパーミッション拒否時に出力されるエラー付きツール実行結果。
Example:
```json
{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","is_error":true,"content":"Error: file not found"}]}}
```

##### result/success

ターンの正常完了を示すメッセージ。Read() はこのメッセージ受信時に読み取りを終了する。
Result フィールドには最後のテキストブロックの内容が入る。
result が空文字列の場合、Result フィールドは省略される（部分一致用）。
Example:
```json
{"type":"result","subtype":"success","result":"Hello!","is_error":false,"duration_ms":1234,"duration_api_ms":1000}
```

##### result/success, is_error=true

max_tokens による応答打ち切り時に出力される。Subtype は "success" だが IsError が true になる。
Example:
```json
{"type":"result","subtype":"success","is_error":true,"result":"truncated response..."}
```

##### result/success, permission_denials

非インタラクティブモードでパーミッション拒否されたツールがある場合に出力される。
Example:
```json
{"type":"result","subtype":"success","permission_denials":[{"tool_name":"AskUserQuestion"}]}
```

##### result/error_during_execution

APIがSSEエラーイベントを返した場合のターン終了メッセージ。
assistant メッセージは出力されず、Errors 配列にエラー種別とメッセージが格納される。
Example:
```json
{"type":"result","subtype":"error_during_execution","errors":[{"type":"overloaded_error","message":"Overloaded"}]}
```

#### MessageBody

MessageBody は assistant/user メッセージの内部構造。

| フィールド | JSON キー | 説明 |
|-----------|-----------|------|
| `Role` | `role` | "assistant" または "user" |
| `Content` | `content` | コンテンツブロックの配列 |

#### ContentBlock

ContentBlock はメッセージ内のコンテンツブロックの統合表現。
[ContentBlockType] に応じて使用されるフィールドが異なる。

| フィールド | JSON キー | 説明 |
|-----------|-----------|------|
| `Type` | `type` | ブロック種別 |
| `Text` | `text` | 応答テキスト (text) |
| `Name` | `name` | ツール名 (tool_use) |
| `ID` | `id` | ツール呼び出しID (tool_use) |
| `Input` | `input` | ツールパラメータ (tool_use) |
| `Thinking` | `thinking` | 思考テキスト (thinking) |
| `ToolUseID` | `tool_use_id` | 対応するツール呼び出しID (tool_result) |
| `Content` | `content` | 実行結果テキスト (tool_result) |
| `IsError` | `is_error` | エラー時 true (tool_result) |

##### type=text

テキスト応答ブロック。モデルの自然言語応答テキストを含む。

JSON例:

```
{"type":"text","text":"The answer is 42."}
```

##### type=tool_use

ツール呼び出しブロック。モデルがツールの実行を要求する際に出力される。

JSON例:

```
{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello","description":"Print hello"}}
```

観測事項:
  - --dangerously-skip-permissions なしの場合、パーミッション拒否で tool_result に is_error:true が返る。

##### type=thinking

拡張思考ブロック。モデルの内部推論過程を含む。

JSON例:

```
{"type":"thinking","thinking":"Let me think about this step by step..."}
```

観測事項:
  - thinking ブロックは text ブロックより前に、別の assistant メッセージとして出力される。

##### type=tool_result

ツール実行結果ブロック。CLIがツール実行結果を報告する。

JSON例:

```
{"type":"tool_result","tool_use_id":"toolu_001","content":"result text","is_error":false}
```

#### ResultError

ResultError は result メッセージ内のエラー詳細。

| フィールド | JSON キー | 説明 |
|-----------|-----------|------|
| `Type` | `type` | エラー種別 (例: "overloaded_error") |
| `Message` | `message` | エラーメッセージ |

#### PermissionDenial

PermissionDenial はパーミッション拒否されたツールの情報。

| フィールド | JSON キー | 説明 |
|-----------|-----------|------|
| `ToolName` | `tool_name` | 拒否されたツール名 |

