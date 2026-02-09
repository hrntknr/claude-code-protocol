# Agent Team

## Team creation via the TeamCreate tool

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a team"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team created.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## TeamDelete tool behavior when no active team exists

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"delete the team"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamDelete","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled team deletion."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Handled team deletion.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## SendMessage tool behavior without team context

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"send a message"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:SendMessage)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"SendMessage","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Handled send message."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Handled send message.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Spawning teammates via the Task tool

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"create team and spawn a teammate"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(tool_use:Task)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"Task","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

## Full agent team lifecycle (create -> delete) across multiple turns

| direction | message | json |
|-----------|---------|------|
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"create a team called proto-test-lifecycle"}}` |
| -> | [system/init](#systeminit) | `{"type":"system","subtype":"init","cwd":"","session_id":"","tools":[],"mcp_servers":[],"model":"","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"","claude_code_version":"","output_style":"","agents":[],"skills":[],"plugins":[],"uuid":""}` |
| -> | [assistant(tool_use:TeamCreate)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamCreate","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team created successfully."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team created successfully.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |
| <- | [user](#user) | `{"type":"user","message":{"role":"user","content":"now delete the team"}}` |
| -> | [assistant(tool_use:TeamDelete)](#assistanttool_use) | `{"type":"assistant","message":{"content":[{"type":"tool_use","id":"","name":"TeamDelete","input":{}}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [user(tool_result)](#usertool_result) | `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"","content":""}]},"session_id":"","uuid":"","tool_use_result":""}` |
| -> | [assistant(text)](#assistanttext) | `{"type":"assistant","message":{"content":[{"type":"text","text":"Team deleted."}],"id":"","model":"","role":"assistant","type":"message","usage":{}},"session_id":"","uuid":""}` |
| -> | [result/success](#resultsuccess) | `{"type":"result","subtype":"success","is_error":false,"duration_ms":0,"duration_api_ms":0,"num_turns":0,"result":"Team deleted.","session_id":"","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":""}` |

