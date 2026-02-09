# Tool Interactive

## AskUserQuestion tool behavior in non-interactive mode

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"ask me a question"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:AskUserQuestion)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"AskUserQuestion","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":"","is_error":true}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"You chose Go. Let me proceed with Go."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"You chose Go. Let me proceed with Go.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[{"tool_name":"AskUserQuestion","tool_use_id":"","tool_input":{}}],"uuid":""}` |

