# Basic

- [Basic flow for a simple text response](#basic-flow-for-a-simple-text-response)
- [Single Bash tool use flow](#single-bash-tool-use-flow)
- [Two-step sequential tool use flow](#two-step-sequential-tool-use-flow)

## Basic flow for a simple text response

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "say hello"
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
        "text": "Hello!"
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
  "result": "Hello!",
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

## Single Bash tool use flow

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "run echo tool-use-test-output"
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
        "text": "The command printed: tool-use-test-output"
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
  "result": "The command printed: tool-use-test-output",
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

## Two-step sequential tool use flow

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "run two echo commands"
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
        "text": "Both commands completed successfully."
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
  "result": "Both commands completed successfully.",
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

