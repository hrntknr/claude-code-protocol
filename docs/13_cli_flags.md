# Cli Flags

- [Replay user messages via --replay-user-messages flag](#replay-user-messages-via---replay-user-messages-flag)
- [Partial message streaming via --include-partial-messages flag](#partial-message-streaming-via---include-partial-messages-flag)
- [Turn limit behavior via --max-turns flag](#turn-limit-behavior-via---max-turns-flag)
- [fast_mode_state defaults to "off" without authentication](#fast_mode_state-defaults-to-off-without-authentication)

## Replay user messages via --replay-user-messages flag

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "replay this message"
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
<tr><td>-&gt;</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "replay this message"
  },
  "session_id": "",
  "uuid": "",
  "isReplay": true
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Echoed!"
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
  "result": "Echoed!",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Partial message streaming via --include-partial-messages flag

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "stream something"
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
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Streamed response."
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
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
  "session_id": "",
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#stream_event">stream_event</a></td><td><pre lang="json">
{
  "type": "stream_event",
  "event": {},
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
  "result": "Streamed response.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Turn limit behavior via --max-turns flag

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "run a command"
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
<tr><td>-&gt;</td><td><a href="../README.md#resulterror_max_turns">result/error_max_turns</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "error_max_turns",
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

## fast_mode_state defaults to "off" without authentication

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "hello"
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
<tr><td>-&gt;</td><td><a href="../README.md#resultsuccess">result/success</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "success",
  "is_error": false,
  "duration_ms": 0,
  "duration_api_ms": 0,
  "num_turns": 0,
  "result": "Hello",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

