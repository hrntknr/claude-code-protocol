# Claude Code CLI Protocol

## Scenarios

- [Basic](docs/01_basic.md)
- [Complex Flows](docs/02_complex_flows.md)
- [Thinking](docs/03_thinking.md)
- [Tool File](docs/04_tool_file.md)
- [Tool Todo](docs/05_tool_todo.md)
- [Tool Task](docs/06_tool_task.md)
- [Tool Plan](docs/07_tool_plan.md)
- [Tool Interactive](docs/08_tool_interactive.md)
- [Tool Web](docs/09_tool_web.md)
- [Tool Restriction](docs/10_tool_restriction.md)
- [Agent Team](docs/11_agent_team.md)
- [Cli Flags](docs/12_cli_flags.md)
- [Error](docs/13_error.md)

## Messages

### system/init

The first message output when the CLI starts. Contains initial information such as session ID and version.
Always output as the first message of a session; the tools field contains the available tools list.
Also includes extended information such as mcp_servers, model, slash_commands, agents, skills, and plugins.

```json
{"type":"system","subtype":"init","cwd":"/path","session_id":"abc","tools":["Bash","Read"],"mcp_servers":[],"model":"claude-opus-4-6","permissionMode":"bypassPermissions","slash_commands":[...],"apiKeySource":"none","claude_code_version":"2.1.37","output_style":"default","agents":[...],"skills":[...],"plugins":[],"uuid":"xxx"}
```

### system/status

A message that notifies system state changes.
Output when the permission mode changes, etc.; the permissionMode field indicates the current mode.

```json
{"type":"system","subtype":"status","status":null,"permissionMode":"plan","uuid":"xxx","session_id":"abc"}
```

### assistant

A message containing the model's response. Each content block is output as a separate assistant message.
The content array contains blocks of type text, tool_use, or thinking.

#### assistant(text)

A text response content block. The text field contains the response text.

```json
{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

#### assistant(tool_use)

A tool use content block. The name field contains the tool name and the input field contains the parameters.

```json
{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

#### assistant(thinking)

An extended thinking content block. The thinking field contains the thinking content.

```json
{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think...","signature":""}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
```

### user

A user message. Used for both input (stdin -> CLI) and output (CLI -> stdout).
For input, content is a string; for output (tool execution results), content is a block array.

```json
{"type":"user","message":{"role":"user","content":"say hello"}}
```

#### user(tool_result)

A message where the CLI reports tool execution results. content is an array of tool_result blocks.
Each block has the corresponding tool_use ID in the tool_use_id field.
parent_tool_use_id, session_id, uuid, and tool_use_result are included at the top level.

```json
{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx","tool_use_result":{}}
```

### result/success

A message indicating successful completion of a turn. Indicates that processing of one turn completed normally.
The result field contains the last text block content.
permission_denials is always present; when empty, it is an empty array [].

```json
{"type":"result","subtype":"success","is_error":false,"duration_ms":55,"duration_api_ms":12,"num_turns":1,"result":"Hello!","stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx"}
```

### result/error_during_execution

A turn-ending message when the API returns an error.
In addition to the same common fields as result/success, the errors array contains error message strings.

```json
{"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":52,"duration_api_ms":18,"num_turns":1,"session_id":"abc","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":["error message"]}
```

### result/error_max_turns

A turn-ending message when --max-turns limit is reached.
The subtype is "error_max_turns" and errors is an empty array.

```json
{"type":"result","subtype":"error_max_turns","is_error":false,"duration_ms":184,"duration_api_ms":30,"num_turns":2,"stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":[]}
```

### stream_event

A message emitted when --include-partial-messages is enabled.
Contains SSE stream events as they arrive from the API.
The event field contains the raw SSE event data (message_start,
content_block_start, content_block_delta, content_block_stop,
message_delta, message_stop).

```json
{"type":"stream_event","event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}},"session_id":"abc","parent_tool_use_id":null,"uuid":"xxx"}
```

### user (replay)

A user message echoed back on stdout when --replay-user-messages is enabled.
Contains the same content as the input, plus session_id, uuid, parent_tool_use_id,
and isReplay:true.

```json
{"type":"user","message":{"role":"user","content":"hello"},"session_id":"abc","parent_tool_use_id":null,"uuid":"xxx","isReplay":true}
```

