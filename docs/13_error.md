# Error

## Error handling on tool execution failure

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"read a missing file"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The file does not exist. Let me handle this error."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The file does not exist. Let me handle this error.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Behavior when receiving API-level SSE error events

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"trigger an error"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [result/error_during_execution](#resulterror_during_execution) | `{"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"","errors":[]}` |

