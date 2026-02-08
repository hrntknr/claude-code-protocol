package ccprotocol_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hrntknr/claudecodeprotocol/utils"
)

func TestSimpleTextResponse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello!"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"say hello"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Hello!"}]}}`,
		`{"type":"result", "subtype":"success", "result":"Hello!"}`,
	)
}

// TestToolUseBash verifies the message flow when the assistant uses the Bash tool once.
// Expected: system → assistant (final text only) → result
// Intermediate tool_use/tool_result messages are NOT emitted to stdout.
func TestToolUseBash(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: API tells CLI to run a Bash command
		utils.ToolUseResponse("toolu_bash_001", "Bash", map[string]any{
			"command":     "echo tool-use-test-output",
			"description": "Print test output",
		}),
		// Request 2: After tool execution, API returns final text
		utils.TextResponse("The command printed: tool-use-test-output"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"run echo tool-use-test-output"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"The command printed: tool-use-test-output"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseMultiStep verifies the message flow with two sequential tool uses.
// The API returns tool_use twice before returning the final text.
func TestToolUseMultiStep(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: First Bash tool use
		utils.ToolUseResponse("toolu_bash_001", "Bash", map[string]any{
			"command":     "echo step-one",
			"description": "First step",
		}),
		// Request 2: Second Bash tool use
		utils.ToolUseResponse("toolu_bash_002", "Bash", map[string]any{
			"command":     "echo step-two",
			"description": "Second step",
		}),
		// Request 3: Final text response
		utils.TextResponse("Both commands completed successfully."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"run two echo commands"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Both commands completed successfully."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// ---------------------------------------------------------------------------
// Complex flows
// ---------------------------------------------------------------------------

// NOTE ON INIT REQUESTS:
// The CLI makes 3 internal haiku-model requests on startup before the main
// opus request: (1) quota check, (2) file-change detection, (3) token count.
// Tests that need tools to actually execute must prepend 3 dummy TextResponse
// entries so those init requests don't consume the intended tool responses.
// Tests with a single response (or where only the last matters) work without
// dummies because the stub repeats the last response for extra requests.

// initResponses returns 3 dummy TextResponse entries to absorb the CLI's
// haiku init requests.
func initResponses() [][]utils.SSEEvent {
	return [][]utils.SSEEvent{
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
	}
}

// withInit prepends 3 init-absorbing dummy responses to the given responses.
func withInit(responses ...[]utils.SSEEvent) [][]utils.SSEEvent {
	return append(initResponses(), responses...)
}

// TestTextAndToolUseInSameResponse verifies the CLI when the API returns
// both a text block and a tool_use block in a single response message.
func TestTextAndToolUseInSameResponse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: text + tool_use in one response
		utils.TextAndToolUseResponse(
			"Let me check that.",
			"toolu_combo_001", "Bash", map[string]any{
				"command":     "echo combined-test",
				"description": "Combined test",
			},
		),
		// Request 2: final text after tool execution
		utils.TextResponse("Done. The output was: combined-test"),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"check and run combined"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Done. The output was: combined-test"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestParallelToolUse verifies the CLI when the API returns multiple tool_use blocks
// in a single response (parallel tool calls).
func TestParallelToolUse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: two tool_use blocks in one response
		utils.MultiToolUseResponse(
			utils.ToolCall{
				ID:   "toolu_par_001",
				Name: "Bash",
				Input: map[string]any{
					"command":     "echo parallel-one",
					"description": "First parallel",
				},
			},
			utils.ToolCall{
				ID:   "toolu_par_002",
				Name: "Bash",
				Input: map[string]any{
					"command":     "echo parallel-two",
					"description": "Second parallel",
				},
			},
		),
		// Request 2: final text after both tool executions
		utils.TextResponse("Both commands ran: parallel-one and parallel-two"),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"run two commands in parallel"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Both commands ran: parallel-one and parallel-two"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestMultiTurnConversation verifies that multiple send/read cycles work
// within the same CLI session (multi-turn conversation).
func TestMultiTurnConversation(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		utils.TextResponse("First answer."),
		utils.TextResponse("Second answer."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	// Turn 1
	s.Send(`{"type":"user","message":{"role":"user","content":"first question"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"First answer."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Turn 2
	s.Send(`{"type":"user","message":{"role":"user","content":"second question"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Second answer."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestMaxTokensStopReason verifies CLI behavior when the API response
// is truncated by the max_tokens limit (stop_reason: "max_tokens").
func TestMaxTokensStopReason(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.MaxTokensTextResponse("This response was truncated because it hit the max"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"generate a very long response"}}`)
	// Observed: The CLI retries the max_tokens response multiple times,
	// emitting the truncated text alternating with a synthetic "API Error"
	// message about the max output token limit. Eventually it produces a
	// result with subtype "success" but is_error true.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text"}]}}`,
		`{"type":"result", "subtype":"success", "is_error":true}`,
	)
}

// TestThinkingResponse verifies CLI behavior when the API returns
// an extended-thinking block followed by a text block.
func TestThinkingResponse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ThinkingResponse(
			"Let me think about this step by step...",
			"The answer is 42.",
		),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"what is the answer?"}}`)
	// Observed: Thinking blocks ARE emitted as a separate assistant message
	// with content[0].type="thinking". Then the text block follows as another
	// assistant message. Result contains only the text.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"thinking","thinking":"Let me think about this step by step..."}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"The answer is 42."}]}}`,
		`{"type":"result", "subtype":"success", "result":"The answer is 42."}`,
	)
}

// ---------------------------------------------------------------------------
// Tool coverage
// ---------------------------------------------------------------------------

// TestToolUseRead verifies the Read tool by reading a temporary file.
func TestToolUseRead(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("file-content-for-read-test"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Read the file
		utils.ToolUseResponse("toolu_read_001", "Read", map[string]any{
			"file_path": testFile,
		}),
		// Request 2: Final text
		utils.TextResponse("The file contains: file-content-for-read-test"),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"read the test file"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"The file contains: file-content-for-read-test"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseWrite verifies the Write tool by creating a new file
// and checking that the file was actually created on disk.
func TestToolUseWrite(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "output.txt")

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Write a new file
		utils.ToolUseResponse("toolu_write_001", "Write", map[string]any{
			"file_path": targetFile,
			"content":   "written-by-cli-test",
		}),
		// Request 2: Final text
		utils.TextResponse("File created successfully."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"write a file"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"File created successfully."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Verify the file was actually written to disk.
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("expected file to exist at %s: %v", targetFile, err)
	}
	if string(content) != "written-by-cli-test" {
		t.Errorf("file content = %q, want %q", string(content), "written-by-cli-test")
	}
}

// TestToolUseEdit verifies the Edit tool by reading then editing a file.
// Edit requires a prior Read in the conversation.
func TestToolUseEdit(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "editable.txt")
	if err := os.WriteFile(testFile, []byte("line1\nold-content\nline3\n"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Read the file first (Edit requires prior Read)
		utils.ToolUseResponse("toolu_read_001", "Read", map[string]any{
			"file_path": testFile,
		}),
		// Request 2: Edit the file
		utils.ToolUseResponse("toolu_edit_001", "Edit", map[string]any{
			"file_path":  testFile,
			"old_string": "old-content",
			"new_string": "new-content",
		}),
		// Request 3: Final text
		utils.TextResponse("File edited successfully."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"edit the file"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"File edited successfully."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Verify the file was actually modified.
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("reading edited file: %v", err)
	}
	if !strings.Contains(string(content), "new-content") {
		t.Errorf("expected edited content to contain 'new-content', got: %s", string(content))
	}
	if strings.Contains(string(content), "old-content") {
		t.Errorf("expected 'old-content' to be replaced, got: %s", string(content))
	}
}

// TestToolUseGlob verifies the Glob tool for file pattern matching.
func TestToolUseGlob(t *testing.T) {
	tmpDir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt", "c.log"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(name), 0644); err != nil {
			t.Fatalf("setup: write %s: %v", name, err)
		}
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Glob for .txt files
		utils.ToolUseResponse("toolu_glob_001", "Glob", map[string]any{
			"pattern": filepath.Join(tmpDir, "*.txt"),
		}),
		// Request 2: Final text
		utils.TextResponse("Found 2 text files."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"find txt files"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Found 2 text files."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseGrep verifies the Grep tool for content search.
func TestToolUseGrep(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(tmpDir, "searchable.txt"),
		[]byte("target-pattern-here\nother-line\n"),
		0644,
	); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Grep for pattern
		utils.ToolUseResponse("toolu_grep_001", "Grep", map[string]any{
			"pattern": "target-pattern",
			"path":    tmpDir,
		}),
		// Request 2: Final text
		utils.TextResponse("Found the pattern in searchable.txt."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"search for target-pattern"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Found the pattern in searchable.txt."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseTodoWrite verifies the TodoWrite tool for task list management.
func TestToolUseTodoWrite(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Create a todo list
		utils.ToolUseResponse("toolu_todo_001", "TodoWrite", map[string]any{
			"todos": []any{
				map[string]any{"content": "First task", "status": "in_progress", "activeForm": "Working on first task"},
				map[string]any{"content": "Second task", "status": "pending", "activeForm": "Preparing second task"},
			},
		}),
		// Request 2: Final text
		utils.TextResponse("Created a todo list with 2 items."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"create a todo list"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Created a todo list with 2 items."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestLongToolChain verifies a multi-step tool chain: Read → Edit → Bash.
func TestLongToolChain(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "chain.txt")
	if err := os.WriteFile(testFile, []byte("original-content"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Step 1: Read the file
		utils.ToolUseResponse("toolu_chain_001", "Read", map[string]any{
			"file_path": testFile,
		}),
		// Step 2: Edit the file
		utils.ToolUseResponse("toolu_chain_002", "Edit", map[string]any{
			"file_path":  testFile,
			"old_string": "original-content",
			"new_string": "modified-content",
		}),
		// Step 3: Verify with Bash
		utils.ToolUseResponse("toolu_chain_003", "Bash", map[string]any{
			"command":     "cat " + testFile,
			"description": "Verify file content",
		}),
		// Step 4: Final text
		utils.TextResponse("Chain complete: read, edited, and verified."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"read, edit, and verify the file"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Chain complete: read, edited, and verified."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Verify the file was actually modified through the chain.
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("reading chain file: %v", err)
	}
	if !strings.Contains(string(content), "modified-content") {
		t.Errorf("expected chain file to contain 'modified-content', got: %s", string(content))
	}
}

// ---------------------------------------------------------------------------
// Advanced flows
// ---------------------------------------------------------------------------

// TestThinkingWithToolUse verifies the CLI when the API returns a thinking block
// followed by a tool_use block, then a final text response.
func TestThinkingWithToolUse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: thinking + tool_use
		utils.ThinkingAndToolUseResponse(
			"I need to run a command to check...",
			"toolu_think_001", "Bash", map[string]any{
				"command":     "echo thinking-tool-test",
				"description": "Test after thinking",
			},
		),
		// Request 2: final text
		utils.TextResponse("After thinking and running the command: thinking-tool-test"),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"think and then run a command"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"After thinking and running the command: thinking-tool-test"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestRequestRecording verifies the stub API's request recording capability
// and observes the structure of what the CLI sends back after tool execution.
func TestRequestRecording(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: tool_use
		utils.ToolUseResponse("toolu_rec_001", "Bash", map[string]any{
			"command":     "echo recorded",
			"description": "Record test",
		}),
		// Request 2: final text
		utils.TextResponse("Done."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"run a recorded command"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Done."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Verify that the stub recorded at least 5 requests:
	// 3 haiku init + 1 opus user message (tool_use) + 1 opus follow-up (tool_result).
	reqs := stub.Requests()
	if len(reqs) < 5 {
		t.Fatalf("expected at least 5 recorded requests, got %d", len(reqs))
	}

	// The 5th request (index 4) should contain messages with a tool_result.
	body := reqs[4].Body
	messages, ok := body["messages"]
	if !ok {
		t.Fatal("follow-up request missing 'messages' field")
	}
	msgList, ok := messages.([]any)
	if !ok {
		t.Fatalf("messages is not an array: %T", messages)
	}
	if len(msgList) == 0 {
		t.Fatal("messages array is empty")
	}
	t.Logf("recorded follow-up request messages count: %d", len(msgList))
	for i, m := range msgList {
		t.Logf("  messages[%d]: %v", i, m)
	}
}

// ---------------------------------------------------------------------------
// Additional tool coverage
// ---------------------------------------------------------------------------

// TestToolUseNotebookEdit verifies the NotebookEdit tool by inserting a cell
// into a Jupyter notebook file.
func TestToolUseNotebookEdit(t *testing.T) {
	tmpDir := t.TempDir()
	nbFile := filepath.Join(tmpDir, "test.ipynb")
	nbContent := `{
  "cells": [
    {"cell_type": "code", "execution_count": null, "id": "cell-1", "metadata": {}, "outputs": [], "source": ["print('hello')"]}
  ],
  "metadata": {"kernelspec": {"display_name": "Python 3", "language": "python", "name": "python3"}, "language_info": {"name": "python", "version": "3.10.0"}},
  "nbformat": 4,
  "nbformat_minor": 5
}`
	if err := os.WriteFile(nbFile, []byte(nbContent), 0644); err != nil {
		t.Fatalf("setup: write notebook: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Insert a new cell
		utils.ToolUseResponse("toolu_nb_001", "NotebookEdit", map[string]any{
			"notebook_path": nbFile,
			"cell_id":       "cell-1",
			"cell_type":     "code",
			"new_source":    "print('world')",
			"edit_mode":     "insert",
		}),
		// Request 2: Final text
		utils.TextResponse("Inserted a new cell into the notebook."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"add a cell to the notebook"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Inserted a new cell into the notebook."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Verify the notebook was modified.
	modified, err := os.ReadFile(nbFile)
	if err != nil {
		t.Fatalf("reading notebook: %v", err)
	}
	if !strings.Contains(string(modified), "print('world')") {
		t.Errorf("expected notebook to contain new cell, got: %s", string(modified))
	}
}

// TestToolUseAskUserQuestion verifies the CLI behavior when the API instructs
// it to use AskUserQuestion. This tool requires user interaction, so we
// observe how the CLI handles it in stream-json mode.
func TestToolUseAskUserQuestion(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: AskUserQuestion tool_use
		utils.ToolUseResponse("toolu_ask_001", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Which language do you prefer?",
					"header":      "Language",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "Go", "description": "The Go language"},
						map[string]any{"label": "Rust", "description": "The Rust language"},
					},
				},
			},
		}),
		// Request 2: Final text (after auto-answer or user interaction)
		utils.TextResponse("You chose Go. Let me proceed with Go."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"ask me a question"}}`)
	// Observed: In non-interactive stream-json mode, AskUserQuestion is emitted
	// as an assistant tool_use, then a user tool_result with is_error:true
	// (content "Answer questions?"). The API then returns the final text.
	// The result includes a permission_denials array listing the denied tool.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"AskUserQuestion"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result","is_error":true}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"You chose Go. Let me proceed with Go."}]}}`,
		`{"type":"result", "subtype":"success", "permission_denials":[{"tool_name":"AskUserQuestion"}]}`,
	)
}

// TestToolUseEnterPlanMode verifies the CLI behavior when the API instructs
// it to use EnterPlanMode.
func TestToolUseEnterPlanMode(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: EnterPlanMode tool_use
		utils.ToolUseResponse("toolu_plan_001", "EnterPlanMode", map[string]any{}),
		// Request 2: Final text after plan mode transition
		utils.TextResponse("I have entered plan mode. Let me explore the codebase."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"plan the implementation"}}`)
	// Observed: EnterPlanMode emits the tool_use as an assistant message,
	// then a system status message with permissionMode:"plan", then the
	// user tool_result with plan mode instructions, then the final text.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"EnterPlanMode"}]}}`,
		`{"type":"system", "subtype":"status", "permissionMode":"plan"}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"I have entered plan mode. Let me explore the codebase."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseWebFetch verifies the CLI behavior when the API instructs it
// to use WebFetch. A local static page is served by the stub for fetching.
// Note: WebFetch internally makes an additional haiku API call to process
// the fetched content, which consumes extra stub responses.
func TestToolUseWebFetch(t *testing.T) {
	stub := &utils.StubAPIServer{
		StaticPages: map[string]string{
			"/static/test-page": "<html><body><h1>Test Page</h1><p>Content for WebFetch test.</p></body></html>",
		},
		Responses: withInit(
			// Request 1: WebFetch tool_use
			utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
				"url":    "", // placeholder, will be set after Start()
				"prompt": "What does the page say?",
			}),
			// Extra responses for WebFetch's internal haiku processing
			utils.TextResponse("ok"),
			utils.TextResponse("ok"),
			utils.TextResponse("ok"),
			// Final text
			utils.TextResponse("The page contains a heading: Test Page"),
		),
	}
	stub.Start()
	defer stub.Close()

	// Update the WebFetch URL to point to our stub's static page.
	// The ToolUseResponse is already built, so we rebuild with the correct URL.
	targetURL := stub.URL() + "/static/test-page"
	stub.Responses = withInit(
		utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
			"url":    targetURL,
			"prompt": "What does the page say?",
		}),
		// Extra responses for WebFetch's internal haiku processing
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("The page contains a heading: Test Page"),
	)

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"fetch the test page"}}`)
	// Observed: WebFetch upgrades HTTP to HTTPS, causing an SSL error when
	// hitting the plain HTTP stub server. The CLI emits the tool_use, then
	// a user tool_result with is_error:true containing the SSL error.
	// The API then returns the next response as final text.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"WebFetch"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result","is_error":true}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// ---------------------------------------------------------------------------
// Error handling flows
// ---------------------------------------------------------------------------

// TestToolError verifies the CLI behavior when a tool execution fails.
// The Read tool is called with a non-existent file path.
func TestToolError(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: Read a non-existent file
		utils.ToolUseResponse("toolu_err_001", "Read", map[string]any{
			"file_path": "/tmp/this-file-does-not-exist-for-test-12345.txt",
		}),
		// Request 2: The API acknowledges the error and responds
		utils.TextResponse("The file does not exist. Let me handle this error."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"read a missing file"}}`)
	// The CLI should handle the tool error gracefully.
	// The API receives the error as a tool_result with is_error=true,
	// then returns a normal text response.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"The file does not exist. Let me handle this error."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestAPIError verifies the CLI behavior when the API returns an SSE error event.
func TestAPIError(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// All requests (including init) get the same error.
		// The CLI should handle the API error.
		utils.ErrorSSEResponse("overloaded_error", "Overloaded"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"trigger an error"}}`)
	// Observed: When the API returns an SSE error event, the CLI emits a
	// result with subtype "error_during_execution" and an "errors" array
	// containing the error details. No assistant messages are emitted.
	output := s.Read()
	utils.AssertOutput(t, output,
		`{"type":"system", "subtype":"init"}`,
		`{"type":"result", "subtype":"error_during_execution"}`,
	)
}

// ---------------------------------------------------------------------------
// Additional content block patterns
// ---------------------------------------------------------------------------

// TestMultipleTextBlocks verifies the CLI behavior when the API returns
// a response with multiple text content blocks.
func TestMultipleTextBlocks(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.MultiTextResponse("First paragraph.", "Second paragraph."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"write two paragraphs"}}`)
	// Observed: Each text content block is emitted as a separate assistant
	// message. The result contains only the LAST text block's content.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"First paragraph."}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Second paragraph."}]}}`,
		`{"type":"result", "subtype":"success", "result":"Second paragraph."}`,
	)
}

// ---------------------------------------------------------------------------
// Agent Teams
// ---------------------------------------------------------------------------
// Agent Teams is an experimental multi-agent orchestration feature (added in
// v2.1.32). It allows a lead session to spawn teammate processes that share
// a task list and communicate via file-based inboxes.
// Enabled via CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1.

// agentTeamEnv returns the environment variable to enable agent teams.
func agentTeamEnv() []string {
	return []string{"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1"}
}

// TestToolUseTeamCreate verifies the CLI behavior when the API instructs it
// to use TeamCreate. TeamCreate creates team config and task directories
// at ~/.claude/teams/{name}/ and ~/.claude/tasks/{name}/.
func TestToolUseTeamCreate(t *testing.T) {
	teamName := "proto-test-team-create"

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: TeamCreate
		utils.ToolUseResponse("toolu_tc_001", "TeamCreate", map[string]any{
			"team_name":   teamName,
			"description": "Protocol test team",
		}),
		// Request 2: Final text
		utils.TextResponse("Team created."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), agentTeamEnv())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"create a team"}}`)
	// Observed: TeamCreate emits the tool_use, then a tool_result containing
	// JSON with team_name, team_file_path, and lead_agent_id. The tool_result
	// is NOT an error (is_error is absent). Then final text and result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Team created."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Clean up team files if created.
	home, _ := os.UserHomeDir()
	os.RemoveAll(filepath.Join(home, ".claude", "teams", teamName))
	os.RemoveAll(filepath.Join(home, ".claude", "tasks", teamName))
}

// TestToolUseTeamDelete verifies the CLI behavior when the API instructs it
// to use TeamDelete. Without an active team, this should produce an error
// in the tool_result.
func TestToolUseTeamDelete(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: TeamDelete (no active team)
		utils.ToolUseResponse("toolu_td_001", "TeamDelete", map[string]any{}),
		// Request 2: Final text
		utils.TextResponse("Handled team deletion."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), agentTeamEnv())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"delete the team"}}`)
	// Observed: TeamDelete without an active team does NOT error. It returns
	// a tool_result with success:true and message "No team name found, nothing
	// to clean up". Then final text and result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"TeamDelete"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Handled team deletion."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseSendMessage verifies the CLI behavior when the API instructs it
// to use SendMessage. Without a team context, this should produce an error.
func TestToolUseSendMessage(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: SendMessage (no team context)
		utils.ToolUseResponse("toolu_sm_001", "SendMessage", map[string]any{
			"type":      "message",
			"recipient": "nonexistent-agent",
			"content":   "Hello from test",
			"summary":   "Test message",
		}),
		// Request 2: Final text
		utils.TextResponse("Handled send message."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), agentTeamEnv())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"send a message"}}`)
	// Observed: SendMessage even without a team context does NOT error.
	// It returns a tool_result with success:true containing routing info
	// (sender: "team-lead", target: "@nonexistent-agent"). The message is
	// written to a file-based inbox regardless. Then final text and result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"SendMessage"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Handled send message."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseTaskSpawnTeammate verifies the CLI behavior when the API instructs
// it to use the Task tool with team_name parameter to spawn a teammate.
// The teammate is a separate CLI process that also hits the stub API.
func TestToolUseTaskSpawnTeammate(t *testing.T) {
	teamName := "proto-test-task-teammate"

	stub := &utils.StubAPIServer{Responses: withInit(
		// Request 1: TeamCreate first (needed for teammate spawn)
		utils.ToolUseResponse("toolu_tc_001", "TeamCreate", map[string]any{
			"team_name":   teamName,
			"description": "Team for Task spawn test",
		}),
		// Request 2: Task to spawn a teammate
		utils.ToolUseResponse("toolu_task_001", "Task", map[string]any{
			"description":   "Test teammate",
			"prompt":        "Say hello and finish",
			"subagent_type": "general-purpose",
			"team_name":     teamName,
			"name":          "worker-1",
		}),
		// Extra responses for the teammate's init + main requests
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("Hello from teammate."),
		// Final text from lead
		utils.TextResponse("Teammate completed its task."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), agentTeamEnv())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"create team and spawn a teammate"}}`)
	// Observed: TeamCreate tool_result → Task tool_use → Task tool_result.
	// The Task tool_result contains status "teammate_spawned" with agent details
	// including agent_id, name, team_name, color, model. The teammate is spawned
	// as a background process (in-process mode). Then final text and result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"Task"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Clean up team files.
	home, _ := os.UserHomeDir()
	os.RemoveAll(filepath.Join(home, ".claude", "teams", teamName))
	os.RemoveAll(filepath.Join(home, ".claude", "tasks", teamName))
}

// TestAgentTeamLifecycle verifies a full agent team lifecycle in a multi-turn
// session: TeamCreate → SendMessage (broadcast) → TeamDelete.
func TestAgentTeamLifecycle(t *testing.T) {
	teamName := "proto-test-lifecycle"

	stub := &utils.StubAPIServer{Responses: withInit(
		// Turn 1, Request 1: TeamCreate
		utils.ToolUseResponse("toolu_tc_001", "TeamCreate", map[string]any{
			"team_name":   teamName,
			"description": "Lifecycle test team",
		}),
		// Turn 1, Request 2: Final text
		utils.TextResponse("Team created successfully."),
		// Turn 2 responses (after user sends second message):
		// The CLI makes additional requests for the second turn.
		// Request 1: TeamDelete
		utils.ToolUseResponse("toolu_td_001", "TeamDelete", map[string]any{}),
		// Request 2: Final text
		utils.TextResponse("Team deleted."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), agentTeamEnv())
	defer s.Close()

	// Turn 1: Create team
	s.Send(`{"type":"user","message":{"role":"user","content":"create a team called ` + teamName + `"}}`)
	// Observed: TeamCreate emits tool_use → tool_result → final text → result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"TeamCreate"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Team created successfully."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Turn 2: Delete team
	s.Send(`{"type":"user","message":{"role":"user","content":"now delete the team"}}`)
	// Observed: TeamDelete in second turn emits init again (CLIのsession状態のリフレッシュ),
	// then tool_use → tool_result with success:true and cleanup message → final text → result.
	utils.AssertOutput(t, s.Read(),
		`{"type":"assistant", "message":{"content":[{"type":"tool_use","name":"TeamDelete"}]}}`,
		`{"type":"user", "message":{"content":[{"type":"tool_result"}]}}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Team deleted."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)

	// Clean up in case TeamDelete didn't work.
	home, _ := os.UserHomeDir()
	os.RemoveAll(filepath.Join(home, ".claude", "teams", teamName))
	os.RemoveAll(filepath.Join(home, ".claude", "tasks", teamName))
}
