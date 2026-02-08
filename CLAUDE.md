A project that analyzes the `--output-format stream-json` protocol of Claude Code CLI through observation and records the results as regression tests.

# Purpose

- The analysis target is solely the stream-json I/O of Claude Code CLI (the `claude` command).
- The architecture is: Test code → Claude Code CLI → Stub API, observing the CLI's input/output.
  - input: JSONL sent to stdin via `--input-format stream-json`
  - output: JSONL received from stdout via `--output-format stream-json`
- Observation results are recorded as Go tests to serve as regression tests that detect breaking changes in the CLI.

# Structure

```
protocol_test.go     -- Main observation/regression test file
utils/
  stub_api.go        -- SSE streaming stub for the Anthropic Messages API
  harness.go         -- CLI process management and assertion utilities
```

# Writing Tests

## Basic Pattern

1. Configure the SSE response sequence to return from `utils.StubAPIServer`, then call `Start()`
2. Launch the CLI process with `utils.NewSession(t, stub.URL())`
3. Send a user message in stream-json format to stdin via `s.Send(...)`
4. Read from stdout up to the `result` message via `s.Read()`
5. Assert with partial matching via `utils.AssertOutput(t, output, patterns...)`

## Response Helpers

- `utils.TextResponse(text)` — Generates SSE event sequence for a text response
- `utils.ToolUseResponse(toolID, toolName, input)` — Generates SSE event sequence for a tool call (stop_reason: "tool_use")

For tool-use scenarios, set multiple responses in `Responses`. Responses are consumed in request order; when exhausted, the last response is repeated.

## AssertOutput

Patterns use JSON partial matching. A pattern matches if the specified fields are contained in the actual message. Patterns are searched in order within the output, and non-matching messages are skipped.

# Assumptions & Constraints

- All input to the CLI must go through stdin; all output must come through stdout.
- The CLI is launched with `--dangerously-skip-permissions --verbose --no-session-persistence`.
- The `ANTHROPIC_BASE_URL` environment variable points to the stub API.
- Do not hard-code observed formats based on speculation. Ignore unknown fields and use partial matching for assertions.
- The implementation language is Go. `protocol_test.go` is the main deliverable.
