# Complex Flows

## Text and tool use in the same response

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"check and run combined"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Done. The output was: combined-test"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Done. The output was: combined-test","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Parallel invocation of multiple tools

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"run two commands in parallel"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Both commands ran: parallel-one and parallel-two"}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Both commands ran: parallel-one and parallel-two","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Multi-turn conversation within the same session

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"first question"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First answer."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"First answer.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"second question"}}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second answer."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Second answer.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Behavior when response is truncated by max_tokens

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"generate a very long response"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":""}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":true,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","stop_reason":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Output format for responses containing multiple text blocks

| direction | message | json |
|-----------|---------|------|
| <- | [user](../README.md#user) | `{"type":"user","message":{"role":"user","content":"write two paragraphs"}}` |
| -> | [system/init](../README.md#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"First paragraph."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [assistant(text)](../README.md#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Second paragraph."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](../README.md#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Second paragraph.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

