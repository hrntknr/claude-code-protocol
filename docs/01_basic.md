# Basic

## Basic flow for a simple text response

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"say hello"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Hello!","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Single Bash tool use flow

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"run echo tool-use-test-output"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"The command printed: tool-use-test-output"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"The command printed: tool-use-test-output","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Two-step sequential tool use flow

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"run two echo commands"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands completed successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Both commands completed successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

