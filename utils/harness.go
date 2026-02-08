package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Session manages an interactive CLI process for multi-turn testing.
type Session struct {
	t     *testing.T
	cmd   *exec.Cmd
	stdin interface {
		Write([]byte) (int, error)
		Close() error
	}
	scanner *bufio.Scanner
	stderr  *strings.Builder
}

// NewSession starts a Claude Code CLI process connected to the given stub API.
func NewSession(t *testing.T, baseURL string) *Session {
	t.Helper()
	return NewSessionWithEnv(t, baseURL, nil)
}

// NewSessionWithEnv starts a Claude Code CLI process with additional environment variables.
// extraEnv is a list of "KEY=VALUE" strings appended to the process environment.
func NewSessionWithEnv(t *testing.T, baseURL string, extraEnv []string) *Session {
	t.Helper()

	cmd := exec.Command("claude",
		"--input-format", "stream-json",
		"--output-format", "stream-json",
		"--dangerously-skip-permissions",
		"--verbose",
		"--no-session-persistence",
	)
	cmd.Env = append(cmd.Environ(),
		"ANTHROPIC_BASE_URL="+baseURL,
	)
	cmd.Env = append(cmd.Env, extraEnv...)
	// Create a new process group so we can kill the entire group
	// (including any teammate subprocesses) during cleanup.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}

	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	return &Session{
		t:       t,
		cmd:     cmd,
		stdin:   stdin,
		scanner: scanner,
		stderr:  &stderrBuf,
	}
}

// Send writes input lines to the CLI's stdin.
func (s *Session) Send(lines ...string) {
	s.t.Helper()
	for _, line := range lines {
		if _, err := s.stdin.Write([]byte(line + "\n")); err != nil {
			s.t.Fatalf("write stdin: %v", err)
		}
	}
}

// Read reads output lines from stdout until a "result" message is received.
func (s *Session) Read() []json.RawMessage {
	s.t.Helper()
	var output []json.RawMessage
	for s.scanner.Scan() {
		line := s.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		cp := make([]byte, len(line))
		copy(cp, line)
		msg := json.RawMessage(cp)
		output = append(output, msg)
		s.t.Logf("output[%d]: %s", len(output)-1, string(msg))

		if extractType(msg) == "result" {
			break
		}
	}
	if err := s.scanner.Err(); err != nil {
		s.t.Fatalf("scan stdout: %v", err)
	}
	if len(output) == 0 {
		s.t.Fatal("no output received from CLI")
	}
	return output
}

// Close closes stdin and waits for the CLI process to exit.
// If the process does not exit within 10 seconds (e.g. due to running
// teammate subprocesses), it is killed via SIGKILL to the process group.
func (s *Session) Close() {
	s.stdin.Close()

	done := make(chan error, 1)
	go func() { done <- s.cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			s.t.Logf("CLI exit: %v (stderr: %s)", err, s.stderr.String())
		}
	case <-time.After(10 * time.Second):
		s.t.Logf("CLI did not exit within 10s, killing process group")
		// Kill the entire process group to clean up child processes (teammates).
		_ = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
		<-done
	}
}

func extractType(msg json.RawMessage) string {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(msg, &m); err != nil {
		return ""
	}
	v, ok := m["type"]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return ""
	}
	return s
}

// MustJSON は v をJSON文字列に変換する。
// テストのアサーションパターン構築に使用する。
func MustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic("MustJSON: " + err.Error())
	}
	return string(b)
}

// AssertOutput verifies that the output contains messages matching each expected
// pattern in order. Each pattern is a JSON string produced by MustJSON.
// The comparison is exact: all fields and array elements must match.
// Messages not matching any pattern are skipped.
// Any マッチャーセンチネルは型レベルの一致で判定される。
// actual 側に存在し pattern 側に存在しないキーは、値が null であれば許容される。
func AssertOutput(t *testing.T, output []json.RawMessage, expectedPatterns ...string) {
	t.Helper()

	pos := 0
	for i, pattern := range expectedPatterns {
		var expect any
		if err := json.Unmarshal([]byte(pattern), &expect); err != nil {
			t.Fatalf("invalid expected pattern [%d]: %v", i, err)
		}

		found := false
		startPos := pos
		for ; pos < len(output); pos++ {
			var actual any
			if err := json.Unmarshal(output[pos], &actual); err != nil {
				continue
			}
			if jsonExact(actual, expect) {
				pos++
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected[%d] not found in output[%d:]:\n  pattern: %s", i, startPos, pattern)
		}
	}
}

// jsonExact returns true if actual matches expect.
// All keys in expect must exist in actual with matching values.
// Extra keys in actual are allowed only if their value is null (nil).
// All elements in arrays must match (same length).
//
// センチネル値の検出:
//   - 文字列 "<any>"          → 任意の値にマッチ（型を問わない）
//   - float64 -1              → 任意の float64 にマッチ
//   - map に "<any>" キー     → 任意の map にマッチ
//   - 配列 ["<any>"]          → 任意の配列にマッチ
func jsonExact(actual, expect any) bool {
	// String sentinel: "<any>" matches any value regardless of type.
	if s, ok := expect.(string); ok && s == "<any>" {
		return true
	}

	switch e := expect.(type) {
	case map[string]any:
		// Map sentinel: {"<any>": true} matches any map.
		if _, ok := e["<any>"]; ok && len(e) == 1 {
			_, ok := actual.(map[string]any)
			return ok
		}
		// Normal map comparison.
		a, ok := actual.(map[string]any)
		if !ok {
			return false
		}
		// All keys in expect must exist in actual with matching values.
		for k, ev := range e {
			av, ok := a[k]
			if !ok {
				return false
			}
			if !jsonExact(av, ev) {
				return false
			}
		}
		// Extra keys in actual are allowed only if their value is null.
		for k, av := range a {
			if _, ok := e[k]; !ok {
				if av != nil {
					return false
				}
			}
		}
		return true
	case []any:
		// Array sentinel: ["<any>"] matches any array.
		if len(e) == 1 {
			if s, ok := e[0].(string); ok && s == "<any>" {
				_, ok := actual.([]any)
				return ok
			}
		}
		// Normal exact array comparison.
		a, ok := actual.([]any)
		if !ok {
			return false
		}
		if len(a) != len(e) {
			return false
		}
		for i, ev := range e {
			if !jsonExact(a[i], ev) {
				return false
			}
		}
		return true
	case float64:
		// Number sentinel: -1 matches any number.
		if e == -1 {
			_, ok := actual.(float64)
			return ok
		}
		a, ok := actual.(float64)
		return ok && a == e
	default:
		return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expect)
	}
}
