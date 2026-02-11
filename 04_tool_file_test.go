package ccprotocol_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// File reading via the Read tool
func TestToolUseRead(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("file-content-for-read-test"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Read the file
		utils.ToolUseResponse("toolu_read_001", "Read", map[string]any{
			"file_path": testFile,
		}),
		// Request 2: Final text
		utils.TextResponse("The file contains: file-content-for-read-test"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "read the test file"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "The file contains: file-content-for-read-test",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "The file contains: file-content-for-read-test"
		}),
	)
}

// File creation via the Write tool
func TestToolUseWrite(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "output.txt")

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Write a new file
		utils.ToolUseResponse("toolu_write_001", "Write", map[string]any{
			"file_path": targetFile,
			"content":   "written-by-cli-test",
		}),
		// Request 2: Final text
		utils.TextResponse("File created successfully."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "write a file"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "File created successfully.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "File created successfully."
		}),
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

// File editing via the Edit tool
func TestToolUseEdit(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "editable.txt")
	if err := os.WriteFile(testFile, []byte("line1\nold-content\nline3\n"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
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
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "edit the file"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "File edited successfully.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "File edited successfully."
		}),
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

// File pattern matching via the Glob tool
func TestToolUseGlob(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt", "c.log"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(name), 0644); err != nil {
			t.Fatalf("setup: write %s: %v", name, err)
		}
	}

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Glob for .txt files
		utils.ToolUseResponse("toolu_glob_001", "Glob", map[string]any{
			"pattern": filepath.Join(tmpDir, "*.txt"),
		}),
		// Request 2: Final text
		utils.TextResponse("Found 2 text files."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "find txt files"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Found 2 text files.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Found 2 text files."
		}),
	)
}

// Content search via the Grep tool
func TestToolUseGrep(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(tmpDir, "searchable.txt"),
		[]byte("target-pattern-here\nother-line\n"),
		0644,
	); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Grep for pattern
		utils.ToolUseResponse("toolu_grep_001", "Grep", map[string]any{
			"pattern": "target-pattern",
			"path":    tmpDir,
		}),
		// Request 2: Final text
		utils.TextResponse("Found the pattern in searchable.txt."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "search for target-pattern"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Found the pattern in searchable.txt.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Found the pattern in searchable.txt."
		}),
	)
}

// Jupyter notebook editing via the NotebookEdit tool
func TestToolUseNotebookEdit(t *testing.T) {
	t.Parallel()
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

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
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
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "add a cell to the notebook"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Inserted a new cell into the notebook.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Inserted a new cell into the notebook."
		}),
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

// Multi-step tool chain: Read -> Edit -> Bash
func TestLongToolChain(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "chain.txt")
	if err := os.WriteFile(testFile, []byte("original-content"), 0644); err != nil {
		t.Fatalf("setup: write test file: %v", err)
	}

	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
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
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "read, edit, and verify the file"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Chain complete: read, edited, and verified.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Chain complete: read, edited, and verified."
		}),
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
