# Tool Web

## URL fetching behavior of the WebFetch tool

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"fetch the test page"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:WebFetch)](../README.md#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"WebFetch","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](../README.md#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":"","is_error":true}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## WebSearch tool behavior (may fail due to test environment restrictions)

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"search the web"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:WebSearch)](../README.md#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"WebSearch","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](../README.md#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

