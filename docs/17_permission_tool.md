# Permission Tool

- [Bash tool permission approved via --permission-prompt-tool stdio.](#bash-tool-permission-approved-via---permission-prompt-tool-stdio)
- [Bash tool permission denied via --permission-prompt-tool stdio.](#bash-tool-permission-denied-via---permission-prompt-tool-stdio)

## Bash tool permission approved via --permission-prompt-tool stdio.

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "remove the test file"
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
  "permissionMode": "default",
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:Bash)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "",
        "name": "Bash",
        "input": {}
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
<tr><td>-&gt;</td><td><a href="../README.md#control_request">control_request</a></td><td><pre lang="json">
{
  "type": "control_request",
  "request_id": "",
  "request": {
    "subtype": "can_use_tool",
    "tool_name": "Bash",
    "input": {},
    "tool_use_id": "",
    "permission_suggestions": [],
    "blocked_path": ""
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
      "updatedInput": {
        "command": "rm -f /tmp/ccprotocol_perm_test_file",
        "description": "Remove test file"
      }
    }
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
        "text": "Command executed successfully."
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
  "result": "Command executed successfully.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [],
  "uuid": ""
}
</pre></td></tr>
</table>

## Bash tool permission denied via --permission-prompt-tool stdio.

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "run rm -rf /"
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
  "permissionMode": "default",
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
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:Bash)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "",
        "name": "Bash",
        "input": {}
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
<tr><td>-&gt;</td><td><a href="../README.md#control_request">control_request</a></td><td><pre lang="json">
{
  "type": "control_request",
  "request_id": "",
  "request": {
    "subtype": "can_use_tool",
    "tool_name": "Bash",
    "input": {},
    "tool_use_id": "",
    "permission_suggestions": [],
    "decision_reason": ""
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
      "behavior": "deny",
      "message": "Denied by test"
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
        "tool_use_id": "",
        "content": "",
        "is_error": true
      }
    ]
  },
  "session_id": "",
  "uuid": "",
  "tool_use_result": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttext">assistant(text)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "text",
        "text": "I understand, I will not run that command."
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
  "result": "I understand, I will not run that command.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [
    {
      "tool_name": "Bash",
      "tool_use_id": "",
      "tool_input": {}
    }
  ],
  "uuid": ""
}
</pre></td></tr>
</table>

