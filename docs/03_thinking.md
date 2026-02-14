# Thinking

- [Response containing extended thinking blocks](#response-containing-extended-thinking-blocks)
- [Flow where tool use follows a thinking block](#flow-where-tool-use-follows-a-thinking-block)

## Response containing extended thinking blocks

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "what is the answer?"
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#systeminit">system/init</a></td><td><pre lang="json">
{
  "type": "system",
  "subtype": "init",
  "cwd": "/home/user/project",
  "session_id": "session-abc123",
  "tools": [
    "Bash",
    "Read",
    "Write",
    "Edit",
    "Glob",
    "Grep"
  ],
  "mcp_servers": [],
  "model": "claude-sonnet-4-5-20250929",
  "permissionMode": "bypassPermissions",
  "slash_commands": [],
  "apiKeySource": "env_variable",
  "claude_code_version": "2.1.0",
  "output_style": "default",
  "agents": [],
  "skills": [],
  "plugins": [],
  "uuid": "uuid-abc123",
  "fast_mode_state": "off"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistantthinking">assistant(thinking)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "thinking",
        "thinking": "Let me think about this step by step...",
        "signature": ""
      }
    ],
    "id": "msg_stub_001",
    "model": "claude-sonnet-4-5-20250929",
    "role": "assistant",
    "type": "message",
    "usage": {
      "input_tokens": 10,
      "output_tokens": 1
    }
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "The answer is 42."
      }
    ],
    "id": "msg_stub_001",
    "model": "claude-sonnet-4-5-20250929",
    "role": "assistant",
    "type": "message",
    "usage": {
      "input_tokens": 10,
      "output_tokens": 1
    }
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#resultsuccess">result/success</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "success",
  "is_error": false,
  "duration_ms": 100,
  "duration_api_ms": 50,
  "num_turns": 1,
  "result": "The answer is 42.",
  "session_id": "session-abc123",
  "total_cost_usd": 0.001,
  "usage": {
    "input_tokens": 10,
    "output_tokens": 1
  },
  "modelUsage": {
    "claude-sonnet-4-5-20250929": {
      "input_tokens": 10,
      "output_tokens": 1
    }
  },
  "permission_denials": [],
  "uuid": "uuid-abc123"
}
</pre></td></tr>
</table>

## Flow where tool use follows a thinking block

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "think and then run a command"
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#systeminit">system/init</a></td><td><pre lang="json">
{
  "type": "system",
  "subtype": "init",
  "cwd": "/home/user/project",
  "session_id": "session-abc123",
  "tools": [
    "Bash",
    "Read",
    "Write",
    "Edit",
    "Glob",
    "Grep"
  ],
  "mcp_servers": [],
  "model": "claude-sonnet-4-5-20250929",
  "permissionMode": "bypassPermissions",
  "slash_commands": [],
  "apiKeySource": "env_variable",
  "claude_code_version": "2.1.0",
  "output_style": "default",
  "agents": [],
  "skills": [],
  "plugins": [],
  "uuid": "uuid-abc123",
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
        "text": "After thinking and running the command: thinking-tool-test"
      }
    ],
    "id": "msg_stub_001",
    "model": "claude-sonnet-4-5-20250929",
    "role": "assistant",
    "type": "message",
    "usage": {
      "input_tokens": 10,
      "output_tokens": 1
    }
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#resultsuccess">result/success</a></td><td><pre lang="json">
{
  "type": "result",
  "subtype": "success",
  "is_error": false,
  "duration_ms": 100,
  "duration_api_ms": 50,
  "num_turns": 1,
  "result": "After thinking and running the command: thinking-tool-test",
  "session_id": "session-abc123",
  "total_cost_usd": 0.001,
  "usage": {
    "input_tokens": 10,
    "output_tokens": 1
  },
  "modelUsage": {
    "claude-sonnet-4-5-20250929": {
      "input_tokens": 10,
      "output_tokens": 1
    }
  },
  "permission_denials": [],
  "uuid": "uuid-abc123"
}
</pre></td></tr>
</table>

