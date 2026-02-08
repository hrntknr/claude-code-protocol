# Claude Code CLI stream-json プロトコル

## シナリオ

### シンプルなテキスト応答の基本フロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"say hello"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Hello!","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 単一のBashツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"run echo tool-use-test-output"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The command printed: tool-use-test-output"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The command printed: tool-use-test-output","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 2段階の連続ツール呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"run two echo commands"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands completed successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Both commands completed successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### テキストとツール呼び出しが同一レスポンスに含まれるケース

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"check and run combined"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Done. The output was: combined-test"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Done. The output was: combined-test","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 複数ツールの並列呼び出しフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"run two commands in parallel"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands ran: parallel-one and parallel-two"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Both commands ran: parallel-one and parallel-two","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 同一セッション内でのマルチターン会話

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"first question"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First answer."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"First answer.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"second question"}}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second answer."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Second answer.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### max_tokensで応答が打ち切られた場合の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"generate a very long response"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":""}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":true,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","stop_reason":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 拡張思考（extended thinking）ブロックを含む応答

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"what is the answer?"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(thinking)](#assistantthinking) | `{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think about this step by step...","signature":""}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The answer is 42."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The answer is 42.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Readツールによるファイル読み取り

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"read the test file"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The file contains: file-content-for-read-test"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The file contains: file-content-for-read-test","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Writeツールによるファイル作成

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"write a file"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"File created successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"File created successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Editツールによるファイル編集

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"edit the file"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"File edited successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"File edited successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Globツールによるファイルパターンマッチ

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"find txt files"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Found 2 text files."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Found 2 text files.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Grepツールによるコンテンツ検索

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"search for target-pattern"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Found the pattern in searchable.txt."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Found the pattern in searchable.txt.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### TodoWriteツールによるタスクリスト管理

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a todo list"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Created a todo list with 2 items."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Created a todo list with 2 items.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Read → Edit → Bashの複数ステップツールチェイン

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"read, edit, and verify the file"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Chain complete: read, edited, and verified."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Chain complete: read, edited, and verified.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### 思考ブロック後にツール呼び出しが続くフロー

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"think and then run a command"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"After thinking and running the command: thinking-tool-test"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"After thinking and running the command: thinking-tool-test","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### スタブAPIのリクエスト記録機能

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"run a recorded command"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Done."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Done.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### NotebookEditツールによるJupyterノートブック編集

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"add a cell to the notebook"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Inserted a new cell into the notebook."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Inserted a new cell into the notebook.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### AskUserQuestionツールの非インタラクティブモードでの挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"ask me a question"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:AskUserQuestion)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"AskUserQuestion","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":"","is_error":true}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"You chose Go. Let me proceed with Go."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"You chose Go. Let me proceed with Go.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[{"tool_name":"AskUserQuestion","tool_use_id":"","tool_input":{}}],"uuid":""}` |

### EnterPlanModeツールによるプランモード遷移

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"plan the implementation"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:EnterPlanMode)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"EnterPlanMode","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [system/status](#systemstatus) | `{"type":"system","subtype":"status","permissionMode":"plan","uuid":"","session_id":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"I have entered plan mode. Let me explore the codebase."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"I have entered plan mode. Let me explore the codebase.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### WebFetchツールのURL取得挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"fetch the test page"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:WebFetch)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"WebFetch","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":"","is_error":true}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### ツール実行失敗時のエラーハンドリング

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"read a missing file"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The file does not exist. Let me handle this error."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The file does not exist. Let me handle this error.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### APIレベルのSSEエラーイベント受信時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"trigger an error"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [result/error_during_execution](#resulterror_during_execution) | `{"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"","errors":[]}` |

### 複数テキストブロックを含む応答の出力形式

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"write two paragraphs"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First paragraph."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second paragraph."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Second paragraph.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### TeamCreateツールによるチーム作成

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a team"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team created.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### TeamDeleteツールのアクティブチームなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"delete the team"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamDelete","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled team deletion."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Handled team deletion.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### SendMessageツールのチームコンテキストなし時の挙動

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"send a message"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:SendMessage)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"SendMessage","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled send message."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Handled send message.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### Taskツールによるチームメイトのスポーン（生成）

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"create team and spawn a teammate"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(tool_use:Task)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"Task","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

### エージェントチームのライフサイクル全体（作成→削除）をマルチターンで

| direction | message | json |
|-----------|---------|------|
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a team called proto-test-lifecycle"}}` |
| → | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| → | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team created successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |
| ← | [user](#user) | `{"type":"user","message":{"role":"user","content":"now delete the team"}}` |
| → | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamDelete","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| → | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team deleted."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| → | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team deleted.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## メッセージ

### system/init

CLI起動時に最初に出力されるメッセージ。セッションIDやバージョンなどの初期情報を含む。
必ずセッションの最初のメッセージとして出力され、tools フィールドに利用可能なツール一覧が含まれる。
mcp_servers, model, slash_commands, agents, skills, plugins などの拡張情報も含む。

```json
{"type":"system","subtype":"init","cwd":"/path","session_id":"abc","tools":["Bash","Read"],"mcp_servers":[],"model":"claude-opus-4-6","permissionMode":"bypassPermissions","slash_commands":[...],"apiKeySource":"none","claude_code_version":"2.1.37","output_style":"default","agents":[...],"skills":[...],"plugins":[],"uuid":"xxx"}
```

### system/status

システム状態の変更を通知するメッセージ。
パーミッションモードの変更時などに出力され、permissionMode フィールドで現在のモードが通知される。

```json
{"type":"system","subtype":"status","status":null,"permissionMode":"plan","uuid":"xxx","session_id":"abc"}
```

### assistant

モデルの応答を含むメッセージ。各コンテンツブロックは個別の assistant メッセージとして出力される。
content 配列には text, tool_use, thinking のいずれかのブロックが含まれる。

#### assistant(text)

テキスト応答のコンテンツブロック。text フィールドに応答テキストが入る。

```json
{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

#### assistant(tool_use)

ツール呼び出しのコンテンツブロック。name フィールドにツール名、input フィールドにパラメータが入る。

```json
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

#### assistant(thinking)

拡張思考のコンテンツブロック。thinking フィールドに思考内容が入る。

```json
{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think...","signature":""}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

### user

ユーザメッセージ。入力（stdin→CLI）と出力（CLI→stdout）の両方で使用される。
入力の場合は content が文字列、出力（ツール実行結果）の場合は content がブロック配列となる。

```json
{"type":"user","message":{"role":"user","content":"say hello"}}
```

#### user(tool_result)

CLIがツール実行結果を報告するメッセージ。content は tool_result ブロックの配列。
各ブロックは対応する tool_use の ID を tool_use_id フィールドに持つ。
parent_tool_use_id, session_id, uuid, tool_use_result がトップレベルに含まれる。

```json
{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx","tool_use_result":{}}
```

### result/success

ターンの正常完了を示すメッセージ。1ターンの処理が正常に終了したことを示す。
result フィールドには最後のテキストブロックの内容が入る。
permission_denials は常に存在し、空の場合は空配列 [] となる。

```json
{"type":"result","subtype":"success","is_error":false,"duration_ms":55,"duration_api_ms":12,"num_turns":1,"result":"Hello!","stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx"}
```

### result/error_during_execution

APIがエラーを返した場合のターン終了メッセージ。
result/success と同じ共通フィールドに加え、errors 配列にエラーメッセージ文字列が格納される。

```json
{"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":52,"duration_api_ms":18,"num_turns":1,"session_id":"abc","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":["error message"]}
```

