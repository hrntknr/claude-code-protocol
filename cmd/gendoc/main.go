// cmd/gendoc generates README.md from Go doc comments.
//
// It parses protocol.go for schema type definitions (enum constants and unified
// struct types) and protocol_test.go for scenario descriptions, then produces
// a two-section markdown document:
//
//  1. Scenarios — end-to-end usage examples extracted from test functions
//  2. Messages — enum constants and struct type definitions
//
// Usage:
//
//	go run ./cmd/gendoc
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	root := findProjectRoot()

	msgTypes := parseMessageTypes(filepath.Join(root, "protocol.go"))
	scenarios := parseScenarios(root, filepath.Join(root, "protocol_test.go"))

	var buf strings.Builder
	writeHeader(&buf)
	writeScenarioSection(&buf, scenarios)
	writeMessageSection(&buf, msgTypes)

	outPath := filepath.Join(root, "README.md")
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
// Message type parsing (protocol.go)
// ---------------------------------------------------------------------------

// messageType represents a message type definition with its doc comment.
type messageType struct {
	name    string // e.g. "SystemInitMessage"
	heading string // e.g. "system/init" (from # heading in doc)
	body    string // doc body after the # heading line
}

// parseMessageTypes extracts message types whose doc comments contain a "# heading" line.
func parseMessageTypes(filename string) []messageType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	var result []messageType
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		if gd.Doc == nil {
			continue
		}
		heading, body := extractDocSection(gd.Doc.Text())
		if heading == "" {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			result = append(result, messageType{
				name:    ts.Name.Name,
				heading: heading,
				body:    body,
			})
		}
	}
	return result
}

// extractDocSection extracts the "# ..." heading and the body after it
// from a type's doc comment text.
func extractDocSection(doc string) (heading, body string) {
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
	return
}

// ---------------------------------------------------------------------------
// Scenario parsing (protocol_test.go)
// ---------------------------------------------------------------------------

// scenario represents a test function with its assert patterns extracted from code.
type scenario struct {
	funcName string
	title    string
	turns    []scenarioTurn
}

// scenarioTurn represents a single user input and the corresponding output patterns.
type scenarioTurn struct {
	input   assertPattern
	outputs []assertPattern
}

// assertPattern represents a single MustJSON(...) pattern.
type assertPattern struct {
	label   string // display label, e.g. "system/init", "assistant(tool_use:Bash)"
	heading string // heading anchor target, e.g. "system/init", "assistant"
	json    string // simplified JSON
}

func parseScenarios(root, filename string) []scenario {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Phase 1: collect Go source text of MustJSON arguments.
	type rawScenario struct {
		funcName string
		title    string
		turns    []sourceTurn
	}
	var raws []rawScenario
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		turns := extractSourceTurns(fset, fn)
		if len(turns) == 0 {
			continue
		}
		raws = append(raws, rawScenario{
			funcName: fn.Name.Name,
			title:    titleFromFunc(fn),
			turns:    turns,
		})
	}

	// Phase 2: evaluate all MustJSON expressions via go run.
	var allSources []string
	for _, r := range raws {
		for _, turn := range r.turns {
			allSources = append(allSources, turn.input)
			allSources = append(allSources, turn.outputs...)
		}
	}
	jsons := evalExpressions(root, allSources)

	// Phase 3: derive labels from JSON and assemble final scenarios.
	var scenarios []scenario
	idx := 0
	for _, r := range raws {
		var turns []scenarioTurn
		for _, turn := range r.turns {
			inputJSON := jsons[idx]
			idx++
			inputLabel, inputHeading := labelFromJSON(inputJSON)

			var outputs []assertPattern
			for range turn.outputs {
				rawJSON := jsons[idx]
				label, heading := labelFromJSON(rawJSON)
				outputs = append(outputs, assertPattern{
					label:   label,
					heading: heading,
					json:    simplifyJSON(rawJSON),
				})
				idx++
			}
			turns = append(turns, scenarioTurn{
				input: assertPattern{
					label:   inputLabel,
					heading: inputHeading,
					json:    simplifyJSON(inputJSON),
				},
				outputs: outputs,
			})
		}
		scenarios = append(scenarios, scenario{
			funcName: r.funcName,
			title:    r.title,
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

// sourceTurn pairs the input (s.Send) with output patterns (utils.AssertOutput).
type sourceTurn struct {
	input   string
	outputs []string
}

// extractSourceTurns walks a test function body to find s.Send and utils.AssertOutput
// calls, pairing each Send's MustJSON argument with the following AssertOutput's patterns.
func extractSourceTurns(fset *token.FileSet, fn *ast.FuncDecl) []sourceTurn {
	var pendingInput string
	var turns []sourceTurn
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if isSendCall(call) && len(call.Args) == 1 {
			if src, ok := extractMustJSONSource(fset, call.Args[0]); ok {
				pendingInput = src
			}
			return true
		}
		if isAssertOutputCall(call) && len(call.Args) >= 3 {
			var sources []string
			for _, arg := range call.Args[2:] {
				if src, ok := extractMustJSONSource(fset, arg); ok {
					sources = append(sources, src)
				}
			}
			if len(sources) > 0 && pendingInput != "" {
				turns = append(turns, sourceTurn{
					input:   pendingInput,
					outputs: sources,
				})
				pendingInput = ""
			}
			return true
		}
		return true
	})
	return turns
}

// isSendCall checks if a call expression is s.Send(...).
func isSendCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Send"
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

// extractMustJSONSource returns the Go source text of a utils.MustJSON(...) argument.
func extractMustJSONSource(fset *token.FileSet, expr ast.Expr) (string, bool) {
	outer, ok := expr.(*ast.CallExpr)
	if !ok || len(outer.Args) != 1 {
		return "", false
	}
	sel, ok := outer.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "MustJSON" {
		return "", false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "utils" {
		return "", false
	}
	return exprSource(fset, outer.Args[0]), true
}

// ---------------------------------------------------------------------------
// Label derivation from JSON
// ---------------------------------------------------------------------------

// labelFromJSON derives a display label and heading anchor from the JSON output.
// It extracts type/subtype for the heading, and appends content block detail for the label.
func labelFromJSON(jsonStr string) (label, heading string) {
	var m map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return "unknown", ""
	}

	msgType, _ := m["type"].(string)
	subtype, _ := m["subtype"].(string)

	heading = msgType
	if subtype != "" {
		heading = msgType + "/" + subtype
	}
	label = heading

	// Extract content block type from message.content array.
	if msg, ok := m["message"].(map[string]any); ok {
		if content, ok := msg["content"].([]any); ok && len(content) > 0 {
			if block, ok := content[0].(map[string]any); ok {
				if blockType, ok := block["type"].(string); ok {
					detail := blockType
					if name, _ := block["name"].(string); name != "" && name != "<any>" {
						detail += ":" + name
					}
					label += "(" + detail + ")"
					heading += "(" + blockType + ")"
				}
			}
		}
	}

	return
}

// ---------------------------------------------------------------------------
// Expression evaluation (go run)
// ---------------------------------------------------------------------------

// exprSource returns the Go source text of an AST expression.
func exprSource(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return buf.String()
}

// evalExpressions builds a temporary Go file that evaluates all source expressions
// via utils.MustJSON and returns the JSON strings.
func evalExpressions(root string, sources []string) []string {
	if len(sources) == 0 {
		return nil
	}

	var src strings.Builder
	src.WriteString("package main\n\nimport (\n\t\"fmt\"\n")
	src.WriteString("\t. \"github.com/hrntknr/claudecodeprotocol\"\n")
	src.WriteString("\t\"github.com/hrntknr/claudecodeprotocol/utils\"\n")
	src.WriteString(")\n\n")
	// Suppress unused import errors.
	src.WriteString("var _ = utils.MustJSON\n")
	src.WriteString("var _ MessageBase\n\n")
	src.WriteString("func main() {\n")
	for _, s := range sources {
		src.WriteString("\tfmt.Println(utils.MustJSON(" + s + "))\n")
	}
	src.WriteString("}\n")

	tmpDir, err := os.MkdirTemp(root, ".gendoc-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(tmpFile, []byte(src.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write temp file: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "run", tmpFile)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "eval expressions: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n")
	if len(lines) != len(sources) {
		fmt.Fprintf(os.Stderr, "eval: expected %d results, got %d\n", len(sources), len(lines))
		os.Exit(1)
	}
	return lines
}

// ---------------------------------------------------------------------------
// JSON simplification (sentinel → empty value)
// ---------------------------------------------------------------------------

// simplifyJSON replaces sentinel values in a compact JSON string with empty values.
// Input is always compact JSON from json.Marshal, so string replacement preserves key order.
// json.Marshal escapes <> to \u003c/\u003e, so sentinels use the escaped form.
// Order matters: composite sentinels must be replaced before the string sentinel.
func simplifyJSON(jsonStr string) string {
	s := jsonStr
	s = strings.ReplaceAll(s, `{"\u003cany\u003e":true}`, `{}`)
	s = strings.ReplaceAll(s, `["\u003cany\u003e"]`, `[]`)
	s = strings.ReplaceAll(s, `"\u003cany\u003e"`, `""`)
	s = strings.ReplaceAll(s, `:-1,`, `:0,`)
	s = strings.ReplaceAll(s, `:-1}`, `:0}`)
	return s
}

// ---------------------------------------------------------------------------
// Markdown generation
// ---------------------------------------------------------------------------

func writeHeader(buf *strings.Builder) {
	buf.WriteString("# Claude Code CLI Protocol\n\n")
}

func writeScenarioSection(buf *strings.Builder, scenarios []scenario) {
	buf.WriteString("## Scenarios\n\n")

	for _, sc := range scenarios {
		buf.WriteString("### " + sc.title + "\n\n")
		buf.WriteString("| direction | message | json |\n")
		buf.WriteString("|-----------|---------|------|\n")

		for _, turn := range sc.turns {
			buf.WriteString("| ← | [" + turn.input.label + "](#" + headingToAnchor(turn.input.heading) + ") | `" + turn.input.json + "` |\n")

			for _, p := range turn.outputs {
				buf.WriteString("| → | [" + p.label + "](#" + headingToAnchor(p.heading) + ") | `" + p.json + "` |\n")
			}
		}

		buf.WriteString("\n")
	}
}

func writeMessageSection(buf *strings.Builder, msgTypes []messageType) {
	buf.WriteString("## Messages\n\n")
	for _, mt := range msgTypes {
		buf.WriteString("### " + mt.heading + "\n\n")
		buf.WriteString(godocToMarkdown(mt.body) + "\n\n")
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
