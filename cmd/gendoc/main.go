// cmd/gendoc generates README.md and per-category docs from Go doc comments.
//
// It parses protocol.go for schema type definitions (enum constants and unified
// struct types) and all *_test.go files in the project root for scenario
// descriptions, then produces:
//
//  1. docs/<category>.md — per-category scenario docs (one per test file)
//  2. README.md — index table linking to category docs + Messages section
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
	"sort"
	"strings"
	"unicode"
)

// testFileScenarios groups scenarios parsed from a single test file.
type testFileScenarios struct {
	filename  string     // e.g. "basic_test.go"
	category  string     // e.g. "basic" (derived from filename)
	title     string     // e.g. "Basic" (human-readable heading)
	notes     string     // blockquote notes extracted from file doc comment
	scenarios []scenario // test scenarios extracted from the file
}

func main() {
	root := findProjectRoot()

	msgTypes := parseMessageTypes(filepath.Join(root, "protocol.go"))
	fileScenarios := parseAllTestFiles(root)

	// Create docs/ directory.
	docsDir := filepath.Join(root, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating docs dir: %v\n", err)
		os.Exit(1)
	}

	// Generate per-category docs.
	for _, fs := range fileScenarios {
		var buf strings.Builder
		buf.WriteString("# " + fs.title + "\n\n")
		if fs.notes != "" {
			buf.WriteString(fs.notes + "\n\n")
		}
		writeScenarioIndex(&buf, fs.scenarios)
		writeScenarioSection(&buf, fs.scenarios)
		outPath := filepath.Join(docsDir, fs.category+".md")
		if err := os.WriteFile(outPath, []byte(buf.String()), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", outPath)
	}

	// Generate README.md with index table + Messages section.
	var buf strings.Builder
	writeHeader(&buf)
	writeIndexTable(&buf, fileScenarios)
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
// Multi-file test scanning
// ---------------------------------------------------------------------------

// parseAllTestFiles globs all *_test.go files in the project root, parses
// scenarios from each, and returns sorted results. Files with no test
// functions (e.g. helpers_test.go) are skipped.
func parseAllTestFiles(root string) []testFileScenarios {
	matches, err := filepath.Glob(filepath.Join(root, "*_test.go"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "glob test files: %v\n", err)
		os.Exit(1)
	}
	sort.Strings(matches)

	var result []testFileScenarios
	for _, path := range matches {
		filename := filepath.Base(path)
		scenarios := parseScenarios(root, path)
		if len(scenarios) == 0 {
			continue
		}
		category := strings.TrimSuffix(filename, "_test.go")
		notes := parseFileNotes(path)
		result = append(result, testFileScenarios{
			filename:  filename,
			category:  category,
			title:     categoryToTitle(category),
			notes:     notes,
			scenarios: scenarios,
		})
	}
	return result
}

// parseFileNotes extracts blockquote lines ("> ...") from the file-level doc comment.
func parseFileNotes(path string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return ""
	}
	if f.Doc == nil {
		return ""
	}
	var notes []string
	for _, line := range strings.Split(f.Doc.Text(), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "> ") {
			notes = append(notes, trimmed)
		}
	}
	return strings.Join(notes, "\n")
}

// categoryToTitle converts a snake_case category name to a title-case heading.
// A leading numeric prefix (e.g. "01_") is stripped before conversion.
// e.g. "01_basic" -> "Basic", "11_agent_team" -> "Agent Team"
func categoryToTitle(category string) string {
	// Strip leading numeric prefix like "01_", "13_".
	s := category
	if i := strings.IndexByte(s, '_'); i > 0 {
		allDigits := true
		for _, r := range s[:i] {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			s = s[i+1:]
		}
	}
	words := strings.Split(s, "_")
	for i, w := range words {
		if len(w) > 0 {
			runes := []rune(w)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
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
// Scenario parsing (*_test.go)
// ---------------------------------------------------------------------------

// scenario represents a test function with its assert patterns extracted from code.
type scenario struct {
	funcName string
	title    string
	turns    []scenarioTurn
}

// scenarioTurn represents one or more stdin inputs and the corresponding output patterns.
// inputs may be empty for continuation turns (e.g. Phase 3 after a control_response
// that uses a runtime variable and is therefore invisible to gendoc's AST scanner).
type scenarioTurn struct {
	inputs  []assertPattern
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
			allSources = append(allSources, turn.inputs...)
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
			var inputs []assertPattern
			for range turn.inputs {
				rawJSON := jsons[idx]
				label, heading := labelFromJSON(rawJSON)
				inputs = append(inputs, assertPattern{
					label:   label,
					heading: heading,
					json:    simplifyJSON(rawJSON),
				})
				idx++
			}

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
				inputs:  inputs,
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

// sourceTurn pairs stdin inputs (s.Send) with output patterns (utils.AssertOutput).
// inputs may contain zero or more entries: zero for continuation turns (when
// a preceding Send uses a runtime variable), or multiple when several Sends
// appear before a single AssertOutput (e.g. control_request + user message).
type sourceTurn struct {
	inputs  []string
	outputs []string
}

// extractSourceTurns walks a test function body to find s.Send and utils.AssertOutput
// calls. MustJSON arguments from Send calls are accumulated until an AssertOutput is
// found, at which point a turn is created with all accumulated inputs. If an
// AssertOutput has no preceding Send (e.g. Phase 3 of a permission prompt test),
// a turn with empty inputs is created to preserve the output patterns.
func extractSourceTurns(fset *token.FileSet, fn *ast.FuncDecl) []sourceTurn {
	var pendingInputs []string
	var turns []sourceTurn
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if isSendCall(call) && len(call.Args) == 1 {
			if src, ok := extractMustJSONSource(fset, call.Args[0]); ok {
				pendingInputs = append(pendingInputs, src)
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
			if len(sources) > 0 {
				turns = append(turns, sourceTurn{
					inputs:  pendingInputs,
					outputs: sources,
				})
				pendingInputs = nil
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

// buildEvalSource builds the Go source for evaluating MustJSON expressions.
// extraDecls is inserted before main() to declare stub variables for identifiers
// that are undefined in the temporary evaluation context (e.g. test-local variables).
// helperFuncs contains function source code extracted from helpers_test.go.
func buildEvalSource(sources []string, extraDecls, helperFuncs string) string {
	var src strings.Builder
	src.WriteString("package main\n\nimport (\n\t\"fmt\"\n")
	src.WriteString("\t. \"github.com/hrntknr/claudecodeprotocol\"\n")
	src.WriteString("\t\"github.com/hrntknr/claudecodeprotocol/utils\"\n")
	src.WriteString(")\n\n")
	// Suppress unused import errors.
	src.WriteString("var _ = utils.MustJSON\n")
	src.WriteString("var _ MessageBase\n\n")
	// Override CLI version so all version-gated fields are included in documentation.
	src.WriteString("func init() { utils.TestCLIVersion = \"99.99.99\" }\n\n")
	if helperFuncs != "" {
		src.WriteString(helperFuncs + "\n")
	}
	if extraDecls != "" {
		src.WriteString(extraDecls + "\n")
	}
	src.WriteString("func main() {\n")
	for _, s := range sources {
		src.WriteString("\tfmt.Println(utils.MustJSON(" + s + "))\n")
	}
	src.WriteString("}\n")
	return src.String()
}

// parseUndefinedIdents extracts identifier names from "undefined: X" compile errors.
func parseUndefinedIdents(stderr string) []string {
	var names []string
	seen := map[string]bool{}
	for _, line := range strings.Split(stderr, "\n") {
		if idx := strings.Index(line, "undefined: "); idx >= 0 {
			name := strings.TrimSpace(line[idx+len("undefined: "):])
			if name != "" && !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// extractHelperFuncs reads helpers_test.go and returns the source code of all
// function declarations for inclusion in the gendoc eval program.
func extractHelperFuncs(root string) string {
	path := filepath.Join(root, "helpers_test.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		printer.Fprint(&buf, fset, fn)
		buf.WriteString("\n\n")
	}
	return buf.String()
}

// evalExpressions builds a temporary Go file that evaluates all source expressions
// via utils.MustJSON and returns the JSON strings.
// If compilation fails due to undefined identifiers (e.g. test-local variables like
// reqID), it automatically adds zero-value string declarations and retries.
func evalExpressions(root string, sources []string) []string {
	if len(sources) == 0 {
		return nil
	}

	helperFuncs := extractHelperFuncs(root)

	tmpDir, err := os.MkdirTemp(root, ".gendoc-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, "main.go")

	// First attempt: no extra declarations.
	src := buildEvalSource(sources, "", helperFuncs)
	if err := os.WriteFile(tmpFile, []byte(src), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write temp file: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "run", tmpFile)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// Check for undefined identifier errors and retry with stub declarations.
		undefs := parseUndefinedIdents(stderr.String())
		if len(undefs) > 0 {
			var decls strings.Builder
			for _, name := range undefs {
				decls.WriteString("var " + name + " string\n")
			}
			src = buildEvalSource(sources, decls.String(), helperFuncs)
			if err := os.WriteFile(tmpFile, []byte(src), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write temp file: %v\n", err)
				os.Exit(1)
			}

			stdout.Reset()
			stderr.Reset()
			cmd = exec.Command("go", "run", tmpFile)
			cmd.Dir = root
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "eval expressions (retry): %v\n%s\n", err, stderr.String())
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "eval expressions: %v\n%s\n", err, stderr.String())
			os.Exit(1)
		}
	}

	lines := strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n")
	if len(lines) != len(sources) {
		fmt.Fprintf(os.Stderr, "eval: expected %d results, got %d\n", len(sources), len(lines))
		os.Exit(1)
	}
	return lines
}

// ---------------------------------------------------------------------------
// JSON simplification (sentinel -> empty value)
// ---------------------------------------------------------------------------

// formatJSON pretty-prints a compact JSON string with 2-space indentation,
// preserving the original key order. It scans raw bytes directly instead of
// unmarshal/marshal (which would sort keys).
func formatJSON(src string) string {
	var buf strings.Builder
	indent := 0
	inString := false
	escaped := false
	n := len(src)

	writeIndent := func() {
		buf.WriteByte('\n')
		for i := 0; i < indent; i++ {
			buf.WriteString("  ")
		}
	}

	for i := 0; i < n; i++ {
		c := src[i]
		if escaped {
			buf.WriteByte(c)
			escaped = false
			continue
		}
		if c == '\\' && inString {
			buf.WriteByte(c)
			escaped = true
			continue
		}
		if c == '"' {
			inString = !inString
			buf.WriteByte(c)
			continue
		}
		if inString {
			buf.WriteByte(c)
			continue
		}
		switch c {
		case '{', '[':
			buf.WriteByte(c)
			// Check if empty
			if i+1 < n && (src[i+1] == '}' || src[i+1] == ']') {
				buf.WriteByte(src[i+1])
				i++
			} else {
				indent++
				writeIndent()
			}
		case '}', ']':
			indent--
			writeIndent()
			buf.WriteByte(c)
		case ',':
			buf.WriteByte(c)
			writeIndent()
		case ':':
			buf.WriteString(": ")
		default:
			buf.WriteByte(c)
		}
	}
	return buf.String()
}

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
	buf.WriteString("[![Regression Test](https://github.com/hrntknr/claude-code-protocol/actions/workflows/test.yml/badge.svg)](https://github.com/hrntknr/claude-code-protocol/actions/workflows/test.yml)\n\n")
	buf.WriteString("> **Unofficial** protocol reference reconstructed from analysis of Claude Code's `stream-json` input/output.\n\n")
	buf.WriteString("**Note**\n")
	buf.WriteString("- This documentation is auto-generated from test cases.\n")
	buf.WriteString("- Automated tests (CI) are run against the latest 3 versions of Claude Code CLI.\n\n")
}

// writeIndexTable writes a Scenarios section with links to per-category docs.
func writeIndexTable(buf *strings.Builder, files []testFileScenarios) {
	buf.WriteString("## Scenarios\n\n")
	for _, fs := range files {
		buf.WriteString("- [" + fs.title + "](docs/" + fs.category + ".md)\n")
	}
	buf.WriteString("\n")
}

func writeScenarioIndex(buf *strings.Builder, scenarios []scenario) {
	for _, sc := range scenarios {
		anchor := headingToAnchor(sc.title)
		buf.WriteString("- [" + sc.title + "](#" + anchor + ")\n")
	}
	buf.WriteString("\n")
}

func writeScenarioSection(buf *strings.Builder, scenarios []scenario) {
	for _, sc := range scenarios {
		buf.WriteString("## " + sc.title + "\n\n")
		buf.WriteString("<table>\n")
		buf.WriteString("<tr><th>direction</th><th>message</th><th>json</th></tr>\n")

		for _, turn := range sc.turns {
			for _, p := range turn.inputs {
				writeScenarioRow(buf, "&lt;-", p)
			}
			for _, p := range turn.outputs {
				writeScenarioRow(buf, "-&gt;", p)
			}
		}

		buf.WriteString("</table>\n\n")
	}
}

func writeScenarioRow(buf *strings.Builder, dir string, p assertPattern) {
	buf.WriteString("<tr><td>" + dir + "</td>")
	buf.WriteString("<td><a href=\"../README.md#" + headingToAnchor(p.heading) + "\">" + p.label + "</a></td>")
	buf.WriteString("<td><pre lang=\"json\">\n" + htmlEscape(formatJSON(p.json)) + "\n</pre></td></tr>\n")
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
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
