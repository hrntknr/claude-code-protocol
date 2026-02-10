# Complex Flows

- [Text and tool use in the same response](#text-and-tool-use-in-the-same-response)
- [Parallel invocation of multiple tools](#parallel-invocation-of-multiple-tools)
- [Multi-turn conversation within the same session](#multi-turn-conversation-within-the-same-session)
- [Behavior when response is truncated by max_tokens](#behavior-when-response-is-truncated-by-max_tokens)
- [Output format for responses containing multiple text blocks](#output-format-for-responses-containing-multiple-text-blocks)

## Text and tool use in the same response

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "check and run combined"
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
  "fast_mode_state": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Done. The output was: combined-test"
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
  "result": "Done. The output was: combined-test",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Parallel invocation of multiple tools

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "run two commands in parallel"
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
  "fast_mode_state": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Both commands ran: parallel-one and parallel-two"
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
  "result": "Both commands ran: parallel-one and parallel-two",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Multi-turn conversation within the same session

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "first question"
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
  "fast_mode_state": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "First answer."
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
  "result": "First answer.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "second question"
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Second answer."
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
  "result": "Second answer.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Behavior when response is truncated by max_tokens

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "generate a very long response"
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
  "fast_mode_state": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": ""
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
  "is_error": true,
  "duration_ms": 0,
  "duration_api_ms": 0,
  "num_turns": 0,
  "result": "",
  "stop_reason": "",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Output format for responses containing multiple text blocks

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "write two paragraphs"
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
  "fast_mode_state": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "First paragraph."
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "Second paragraph."
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
  "result": "Second paragraph.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

