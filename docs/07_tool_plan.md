# Tool Plan

## Plan mode transition via the EnterPlanMode tool

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"plan the implementation"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:EnterPlanMode)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"EnterPlanMode","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [system/status](#systemstatus) | `{"type":"system","subtype":"status","permissionMode":"plan","uuid":"","session_id":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"I have entered plan mode. Let me explore the codebase."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"I have entered plan mode. Let me explore the codebase.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## ExitPlanMode tool transitions from plan mode to implementation mode

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"exit plan mode"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:ExitPlanMode)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"ExitPlanMode","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":"","is_error":true}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Plan approved, proceeding with implementation."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Plan approved, proceeding with implementation.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[{"tool_name":"ExitPlanMode","tool_use_id":"","tool_input":{}}],"uuid":""}` |

