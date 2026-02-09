# Cli Flags

## Replay user messages via --replay-user-messages flag

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"replay this message"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"replay this message"},"session_id":"","uuid":"","isReplay":true}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Echoed!"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Echoed!","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Partial message streaming via --include-partial-messages flag

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"stream something"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Streamed response."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [stream_event](../README.md#stream_event) | `{"type":"stream_event","event":{},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Streamed response.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Turn limit behavior via --max-turns flag

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"run a command"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [result/error_max_turns](../README.md#resulterror_max_turns) | `{"type":"result","subtype":"error_max_turns","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"","errors":[]}` |

