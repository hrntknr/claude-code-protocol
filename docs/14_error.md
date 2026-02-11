# Error

- [Error handling on tool execution failure](#error-handling-on-tool-execution-failure)
- [Behavior when receiving API-level SSE error events](#behavior-when-receiving-api-level-sse-error-events)

## Error handling on tool execution failure

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "read a missing file"
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#systeminit">system/init</a></td><td><pre lang="json">
{
  "type": "system",
  "subtype": "init",
  "cwd": "",
  "session_id": "",
  "tools": [],
  "mcp_servers": [],
  "model": "",
  "permissionMode": "bypassPermissions",
  "slash_commands": [],
  "apiKeySource": "",
  "claude_code_version": "",
  "output_style": "",
  "agents": [],
  "skills": [],
  "plugins": [],
  "uuid": "",
  "fast_mode_state": "off"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "The file does not exist. Let me handle this error."
      }
    ],
    "id": "",
    "model": "",
    "role": "assistant",
    "type": "message",
    "usage": {}
  },
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#resultsuccess">result/success</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "success",
  "is_error": false,
  "duration_ms": 0,
  "duration_api_ms": 0,
  "num_turns": 0,
  "result": "The file does not exist. Let me handle this error.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Behavior when receiving API-level SSE error events

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "trigger an error"
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#systeminit">system/init</a></td><td><pre lang="json">
{
  "type": "system",
  "subtype": "init",
  "cwd": "",
  "session_id": "",
  "tools": [],
  "mcp_servers": [],
  "model": "",
  "permissionMode": "bypassPermissions",
  "slash_commands": [],
  "apiKeySource": "",
  "claude_code_version": "",
  "output_style": "",
  "agents": [],
  "skills": [],
  "plugins": [],
  "uuid": "",
  "fast_mode_state": "off"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#resulterror_during_execution">result/error_during_execution</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "error_during_execution",
  "is_error": false,
  "duration_ms": 0,
  "duration_api_ms": 0,
  "num_turns": 0,
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": "",
  "errors": []
}
</pre></td></tr>
</table>

