# Tool Interactive

- [AskUserQuestion tool behavior in non-interactive mode](#askuserquestion-tool-behavior-in-non-interactive-mode)

## AskUserQuestion tool behavior in non-interactive mode

<table>
<tr><th>direction</th><th>message</th><th>json</th></tr>
<tr><td>&lt;-</td><td><a href="../README.md#user">user</a></td><td><pre lang="json">
{
  "type": "user",
  "message": {
    "role": "user",
    "content": "ask me a question"
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
  "uuid": ""
}
</pre></td></tr>
<tr><td>-&gt;</td><td><a href="../README.md#assistanttool_use">assistant(tool_use:AskUserQuestion)</a></td><td><pre lang="json">
{
  "type": "assistant",
  "message": {
    "content": [
      {
        "type": "tool_use",
        "id": "",
        "name": "AskUserQuestion",
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
        "text": "You chose Go. Let me proceed with Go."
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
  "result": "You chose Go. Let me proceed with Go.",
  "session_id": "",
  "total_cost_usd": 0,
  "usage": {},
  "modelUsage": {},
  "permission_denials": [
    {
      "tool_name": "AskUserQuestion",
      "tool_use_id": "",
      "tool_input": {}
    }
  ],
  "uuid": ""
}
</pre></td></tr>
</table>

