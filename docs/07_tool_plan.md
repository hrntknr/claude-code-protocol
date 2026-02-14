# Tool Plan

- [Plan mode transition via the EnterPlanMode tool](#plan-mode-transition-via-the-enterplanmode-tool)
- [ExitPlanMode success via --permission-prompt-tool stdio.](#exitplanmode-success-via---permission-prompt-tool-stdio)
- [ExitPlanMode tool transitions from plan mode to implementation mode](#exitplanmode-tool-transitions-from-plan-mode-to-implementation-mode)

## Plan mode transition via the EnterPlanMode tool

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "plan the implementation"
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:EnterPlanMode)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "toolu_stub_001",
        "name": "EnterPlanMode",
        "input": {
          "command": "echo hello",
          "description": "Example"
        }
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
<tr><td>-&gt;</td><td><a href="../README.md#systemstatus">system/status</a></td><td><pre lang="json">
{
  "type": "system",
  "subtype": "status",
  "permissionMode": "plan",
  "uuid": "uuid-abc123",
  "session_id": "session-abc123"
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#usertool_result">user(tool_result)</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": [
      {
        "type": "tool_result",
        "tool_use_id": "toolu_stub_001",
        "content": "tool execution output"
      }
    ]
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123",
  "tool_use_result": {
    "stdout": "command output"
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
        "text": "I have entered plan mode. Let me explore the codebase."
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
  "result": "I have entered plan mode. Let me explore the codebase.",
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

## ExitPlanMode success via --permission-prompt-tool stdio.

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "exit plan mode"
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
  "permissionMode": "default",
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:ExitPlanMode)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "toolu_stub_001",
        "name": "ExitPlanMode",
        "input": {
          "command": "echo hello",
          "description": "Example"
        }
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
<tr><td>-&gt;</td><td><a href="../README.md#control_request">control_request</a></td><td><pre lang="json">
{
  "type": "control_request",
  "request_id": "request-abc123",
  "request": {
    "subtype": "can_use_tool",
    "tool_name": "ExitPlanMode",
    "input": {
      "command": "echo hello",
      "description": "Example"
    },
    "tool_use_id": "toolu_stub_001"
  }
}
</pre></td></tr>
<tr><td>&lt;-</td><td><a href="../README.md#control_response">control_response</a></td><td><pre lang="json">
{
  "type": "control_response",
  "response": {
    "subtype": "success",
    "request_id": "",
    "response": {
      "behavior": "allow",
      "updatedInput": {}
    }
  }
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#usertool_result">user(tool_result)</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": [
      {
        "type": "tool_result",
        "tool_use_id": "toolu_stub_001",
        "content": "tool execution output"
      }
    ]
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123",
  "tool_use_result": {
    "stdout": "command output"
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
        "text": "Plan approved, proceeding."
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
  "result": "Plan approved, proceeding.",
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

## ExitPlanMode tool transitions from plan mode to implementation mode

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "exit plan mode"
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:ExitPlanMode)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "toolu_stub_001",
        "name": "ExitPlanMode",
        "input": {
          "command": "echo hello",
          "description": "Example"
        }
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
<tr><td>-&gt;</td><td><a href="../README.md#usertool_result">user(tool_result)</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": [
      {
        "type": "tool_result",
        "tool_use_id": "toolu_stub_001",
        "content": "tool execution output",
        "is_error": true
      }
    ]
  },
  "session_id": "session-abc123",
  "uuid": "uuid-abc123",
  "tool_use_result": {
    "stdout": "command output"
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
        "text": "Plan approved, proceeding with implementation."
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
  "result": "Plan approved, proceeding with implementation.",
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
  "permission_denials": [
    {
      "tool_name": "ExitPlanMode",
      "tool_use_id": "toolu_stub_001",
      "tool_input": {
        "key": "value"
      }
    }
  ],
  "uuid": "uuid-abc123"
}
</pre></td></tr>
</table>

