// cmd/gendoc generates PROTOCOL.md from Go doc comments.
//
// It parses protocol.go for schema type definitions (enum constants and unified
// struct types) and protocol_test.go for scenario descriptions, then produces
// a two-section markdown document:
//
//  1. シナリオ — end-to-end usage examples extracted from test functions
//  2. スキーマ — enum constants and struct type definitions
//
// Usage:
//
//	go run ./cmd/gendoc
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ccprotocol "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

func main() {
	root := findProjectRoot()

	docFile := filepath.Join(root, "protocol.go")
	msgFuncs := parseMessageFuncs(docFile)
	scenarios := parseScenarios(filepath.Join(root, "protocol_test.go"))

	var buf strings.Builder
	writeHeader(&buf)
	writeScenarioSection(&buf, scenarios)
	writeMessageSection(&buf, msgFuncs)

	outPath := filepath.Join(root, "PROTOCOL.md")
	if err := os.WriteFile(outPath, []byte(buf.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
		os.Exit(1)
	}
	fmt.Printf("Generated %s\n", outPath)
}

// ---------------------------------------------------------------------------
// Project root detection
// ---------------------------------------------------------------------------

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get working directory:", err)
		os.Exit(1)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Fprintln(os.Stderr, "go.mod not found")
			os.Exit(1)
		}
		dir = parent
	}
}

// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Message function parsing (protocol.go)
// ---------------------------------------------------------------------------

// messageFunc represents a NewMessage* constructor function with its doc comment.
type messageFunc struct {
	name    string // e.g. "NewMessageSystemInit"
	heading string // e.g. "system/init" (from # heading in doc)
	body    string // doc body after the # heading line
}

func parseMessageFuncs(filename string) []messageFunc {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	var result []messageFunc
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !strings.HasPrefix(fn.Name.Name, "NewMessage") {
			continue
		}
		if fn.Doc == nil {
			continue
		}
		heading, body := extractFuncDocSection(fn.Doc.Text())
		result = append(result, messageFunc{
			name:    fn.Name.Name,
			heading: heading,
			body:    body,
		})
	}
	return result
}

// extractFuncDocSection extracts the "# ..." heading and the body after it
// from a function's doc comment text.
func extractFuncDocSection(doc string) (heading, body string) {
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			heading = strings.TrimPrefix(trimmed, "# ")
			if i+1 < len(lines) {
				body = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			}
			return
		}
	}
	// Fallback: no # heading found, use function name and full body after first line.
	heading = ""
	body = skipFirstLine(doc)
	return
}

// ---------------------------------------------------------------------------
// Scenario parsing (protocol_test.go)
// ---------------------------------------------------------------------------

// scenario represents a test function with its assert patterns extracted from code.
type scenario struct {
	funcName string
	title    string
	turns    [][]assertPattern // each turn has ordered assert patterns
}

// assertPattern represents a single MustJSON(NewMessage*(...)) pattern from AssertOutput.
type assertPattern struct {
	label   string // display label, e.g. "system/init", "assistant(tool_use:Bash)"
	heading string // corresponding ##### heading, e.g. "system/init", "assistant/tool_use"
	json    string // actual JSON output from MustJSON
}

func parseScenarios(filename string) []scenario {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	var scenarios []scenario

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		turns := extractTurnsFromFunc(fn)
		if len(turns) == 0 {
			continue
		}
		title := titleFromFunc(fn)
		scenarios = append(scenarios, scenario{
			funcName: fn.Name.Name,
			title:    title,
			turns:    turns,
		})
	}

	return scenarios
}

// titleFromFunc returns a scenario title from the function's doc comment (first line)
// or falls back to the function name (stripping "Test" prefix).
func titleFromFunc(fn *ast.FuncDecl) string {
	if fn.Doc != nil {
		firstLine := strings.TrimSpace(strings.SplitN(fn.Doc.Text(), "\n", 2)[0])
		if firstLine != "" {
			return firstLine
		}
	}
	return fn.Name.Name[len("Test"):]
}

// extractTurnsFromFunc walks a test function body to find utils.AssertOutput calls
// and extracts MustJSON(NewMessage*(...)) patterns from each call.
func extractTurnsFromFunc(fn *ast.FuncDecl) [][]assertPattern {
	var turns [][]assertPattern
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || !isAssertOutputCall(call) {
			return true
		}
		// Args: t, output, patterns...
		if len(call.Args) < 3 {
			return true
		}
		var patterns []assertPattern
		for _, arg := range call.Args[2:] {
			if p, ok := extractAssertPattern(arg); ok {
				patterns = append(patterns, p)
			}
		}
		if len(patterns) > 0 {
			turns = append(turns, patterns)
		}
		return true
	})
	return turns
}

// isAssertOutputCall checks if a call expression is utils.AssertOutput(...).
func isAssertOutputCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "utils" && sel.Sel.Name == "AssertOutput"
}

// extractAssertPattern extracts label and JSON from a utils.MustJSON(NewMessage*(...)) expression.
func extractAssertPattern(expr ast.Expr) (assertPattern, bool) {
	outer, ok := expr.(*ast.CallExpr)
	if !ok || len(outer.Args) != 1 {
		return assertPattern{}, false
	}
	sel, ok := outer.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "MustJSON" {
		return assertPattern{}, false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "utils" {
		return assertPattern{}, false
	}
	inner, ok := outer.Args[0].(*ast.CallExpr)
	if !ok {
		return assertPattern{}, false
	}
	innerIdent, ok := inner.Fun.(*ast.Ident)
	if !ok {
		return assertPattern{}, false
	}
	label, heading, jsonStr := evalConstructor(innerIdent.Name, inner.Args)
	return assertPattern{label: label, heading: heading, json: jsonStr}, true
}

// evalConstructor maps a NewMessage* constructor name and its AST arguments
// to a display label, heading (matching ##### in the generated doc), and JSON string.
func evalConstructor(funcName string, args []ast.Expr) (label, heading, jsonStr string) {
	switch funcName {
	case "NewMessageSystemInit":
		return "system/init", "system/init",
			utils.MustJSON(ccprotocol.NewMessageSystemInit())
	case "NewMessageSystemStatus":
		mode := extractStringLit(args[0])
		return "system/status(" + mode + ")", "system/status",
			utils.MustJSON(ccprotocol.NewMessageSystemStatus(mode))
	case "NewMessageAssistantText":
		text := extractStringLit(args[0])
		return "assistant(text)", "assistant/text",
			utils.MustJSON(ccprotocol.NewMessageAssistantText(text))
	case "NewMessageAssistantToolUse":
		name := extractStringLit(args[0])
		return "assistant(tool_use:" + name + ")", "assistant/tool_use",
			utils.MustJSON(ccprotocol.NewMessageAssistantToolUse(name))
	case "NewMessageAssistantThinking":
		thinking := extractStringLit(args[0])
		return "assistant(thinking)", "assistant/thinking",
			utils.MustJSON(ccprotocol.NewMessageAssistantThinking(thinking))
	case "NewMessageUserToolResult":
		return "user(tool_result)", "user/tool_result",
			utils.MustJSON(ccprotocol.NewMessageUserToolResult())
	case "NewMessageUserToolResultError":
		return "user(tool_result:error)", "user/tool_result, is_error=true",
			utils.MustJSON(ccprotocol.NewMessageUserToolResultError())
	case "NewMessageResultSuccess":
		result := ""
		if len(args) > 0 {
			result = extractStringLit(args[0])
		}
		return "result/success", "result/success",
			utils.MustJSON(ccprotocol.NewMessageResultSuccess(result))
	case "NewMessageResultSuccessIsError":
		return "result/success(is_error:true)", "result/success, is_error=true",
			utils.MustJSON(ccprotocol.NewMessageResultSuccessIsError())
	case "NewMessageResultSuccessWithDenials":
		denials := extractPermissionDenials(args)
		return "result/success(permission_denials)", "result/success, permission_denials",
			utils.MustJSON(ccprotocol.NewMessageResultSuccessWithDenials(denials...))
	case "NewMessageResultErrorDuringExecution":
		return "result/error_during_execution", "result/error_during_execution",
			utils.MustJSON(ccprotocol.NewMessageResultErrorDuringExecution())
	default:
		return funcName, "", "{}"
	}
}

// extractStringLit extracts the string value from a *ast.BasicLit expression.
func extractStringLit(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return ""
	}
	return s
}

// extractPermissionDenials extracts PermissionDenial struct literals from AST arguments.
func extractPermissionDenials(args []ast.Expr) []ccprotocol.PermissionDenial {
	var denials []ccprotocol.PermissionDenial
	for _, arg := range args {
		cl, ok := arg.(*ast.CompositeLit)
		if !ok {
			continue
		}
		var pd ccprotocol.PermissionDenial
		for _, elt := range cl.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if key.Name == "ToolName" {
				pd.ToolName = extractStringLit(kv.Value)
			}
		}
		denials = append(denials, pd)
	}
	return denials
}

// ---------------------------------------------------------------------------
// Markdown generation
// ---------------------------------------------------------------------------

func writeHeader(buf *strings.Builder) {
	buf.WriteString("# Claude Code CLI stream-json プロトコル\n\n")
}

func writeScenarioSection(buf *strings.Builder, scenarios []scenario) {
	buf.WriteString("## シナリオ\n\n")

	for _, sc := range scenarios {
		buf.WriteString("### " + sc.title + "\n\n")
		buf.WriteString("| direction | message | json |\n")
		buf.WriteString("|-----------|---------|------|\n")

		for _, turn := range sc.turns {
			inputJSON := `{"type":"user","message":{"role":"user","content":"..."}}`
			buf.WriteString("| ← | user | `" + inputJSON + "` |\n")

			for _, p := range turn {
				if p.heading != "" {
					buf.WriteString("| → | [" + p.label + "](#" + headingToAnchor(p.heading) + ") | `" + p.json + "` |\n")
				} else {
					buf.WriteString("| → | " + p.label + " | `" + p.json + "` |\n")
				}
			}
		}

		buf.WriteString("\n")
	}
}

func writeMessageSection(buf *strings.Builder, msgFuncs []messageFunc) {
	buf.WriteString("## メッセージ\n\n")
	for _, mf := range msgFuncs {
		buf.WriteString("### " + mf.heading + "\n\n")
		buf.WriteString(godocToMarkdown(mf.body) + "\n\n")
	}
}

// godocToMarkdown converts godoc-formatted text to markdown.
// Explicit fenced code blocks (```json etc.) are passed through as-is,
// with tab-indented content inside them having the tab prefix stripped.
// Tab-indented lines outside fenced blocks are wrapped in auto-generated ``` blocks.
func godocToMarkdown(doc string) string {
	lines := strings.Split(doc, "\n")
	var buf strings.Builder
	inFencedBlock := false // inside explicit ``` block
	inAutoBlock := false   // inside auto-generated ``` block

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		isIndented := strings.HasPrefix(line, "\t")
		isFence := strings.HasPrefix(trimmed, "```")

		if isFence {
			if inAutoBlock {
				buf.WriteString("```\n")
				inAutoBlock = false
			}
			inFencedBlock = !inFencedBlock
			buf.WriteString(trimmed + "\n")
			continue
		}

		if inFencedBlock {
			if isIndented {
				buf.WriteString(strings.TrimPrefix(line, "\t") + "\n")
			} else {
				buf.WriteString(line + "\n")
			}
			continue
		}

		if isIndented && !inAutoBlock {
			buf.WriteString("```\n")
			inAutoBlock = true
		} else if !isIndented && inAutoBlock {
			buf.WriteString("```\n")
			inAutoBlock = false
		}

		if isIndented {
			buf.WriteString(strings.TrimPrefix(line, "\t") + "\n")
		} else {
			buf.WriteString(line + "\n")
		}
	}
	if inAutoBlock {
		buf.WriteString("```\n")
	}
	if inFencedBlock {
		buf.WriteString("```\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}

// headingToAnchor converts a markdown heading to a GitHub-style anchor.
// Lowercase, keep a-z 0-9 _ -, replace spaces with -, remove everything else.
func headingToAnchor(heading string) string {
	heading = strings.ToLower(heading)
	var buf strings.Builder
	for _, r := range heading {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_', r == '-':
			buf.WriteRune(r)
		case r == ' ':
			buf.WriteRune('-')
		}
	}
	return buf.String()
}

// skipFirstLine returns the doc text after the first line, trimmed.
func skipFirstLine(doc string) string {
	if idx := strings.Index(doc, "\n"); idx >= 0 {
		return strings.TrimSpace(doc[idx+1:])
	}
	return ""
}
