# Tool Todo

## Task list management via the TodoWrite tool

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a todo list"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Created a todo list with 2 items."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Created a todo list with 2 items.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

