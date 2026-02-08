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
	"regexp"
	"strconv"
	"strings"

	ccprotocol "github.com/hrntknr/claudecodeprotocol"
)

func main() {
	root := findProjectRoot()

	docFile := filepath.Join(root, "protocol.go")
	enums := parseEnumTypes(docFile)
	structs := parseStructTypes(docFile)
	msgFuncs := parseMessageFuncs(docFile)
	scenarios := parseScenarios(filepath.Join(root, "protocol_test.go"))

	var buf strings.Builder
	writeHeader(&buf)
	writeScenarioSection(&buf, scenarios)
	writeSchemaSection(&buf, enums, structs, msgFuncs)

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
// Enum parsing (protocol.go)
// ---------------------------------------------------------------------------

// enumType represents an enum type with its constants.
type enumType struct {
	name    string        // e.g. "MessageType"
	doc     string        // type doc comment
	consts  []enumConst   // constants in the const block
}

// enumConst represents a single constant in a const block.
type enumConst struct {
	name    string // e.g. "TypeSystem"
	value   string // e.g. "system"
	comment string // doc comment
}

// knownEnumTypes lists the enum type names we want to extract, in order.
var knownEnumTypes = []string{"MessageType", "Subtype", "ContentBlockType"}

func parseEnumTypes(filename string) []enumType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Collect type declarations for enum types (type X string).
	typeDoc := make(map[string]string) // type name -> doc comment
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if isKnownEnum(ts.Name.Name) {
				typeDoc[ts.Name.Name] = cleanDoc(gd.Doc.Text())
			}
		}
	}

	// Collect const blocks. We associate constants with their enum type
	// by looking at the value spec type.
	enumMap := make(map[string][]enumConst) // type name -> consts

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) == 0 {
				continue
			}

			// Determine the enum type from the type annotation.
			typeName := ""
			if vs.Type != nil {
				if ident, ok := vs.Type.(*ast.Ident); ok {
					typeName = ident.Name
				}
			}
			if typeName == "" || !isKnownEnum(typeName) {
				continue
			}

			// Extract the string value.
			value := ""
			if len(vs.Values) > 0 {
				if bl, ok := vs.Values[0].(*ast.BasicLit); ok {
					value = strings.Trim(bl.Value, `"`)
				}
			}

			// Doc comment for the constant.
			comment := cleanDoc(vs.Doc.Text())

			enumMap[typeName] = append(enumMap[typeName], enumConst{
				name:    vs.Names[0].Name,
				value:   value,
				comment: comment,
			})
		}
	}

	// Build result in the order of knownEnumTypes.
	var result []enumType
	for _, name := range knownEnumTypes {
		et := enumType{
			name:   name,
			doc:    typeDoc[name],
			consts: enumMap[name],
		}
		result = append(result, et)
	}
	return result
}

func isKnownEnum(name string) bool {
	for _, n := range knownEnumTypes {
		if n == name {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Struct type parsing (protocol.go)
// ---------------------------------------------------------------------------

// structType represents a struct type definition.
type structType struct {
	name     string
	doc      string // full doc comment
	fields   []structField
	sections []docSection // parsed # type=... sections from doc comment
}

// structField represents a struct field.
type structField struct {
	name    string
	jsonTag string
	comment string
}

// docSection represents a `# type=...` section in a doc comment.
type docSection struct {
	heading string   // e.g. "type=system, subtype=init"
	body    string   // the rest of the section content
}

// knownStructTypes lists the struct type names we want to extract, in order.
var knownStructTypes = []string{"Message", "MessageBody", "ContentBlock", "ResultError", "PermissionDenial"}

func parseStructTypes(filename string) []structType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	structMap := make(map[string]structType)

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if !isKnownStruct(ts.Name.Name) {
				continue
			}

			st := structType{
				name: ts.Name.Name,
				doc:  cleanDoc(gd.Doc.Text()),
			}

			// Parse doc sections.
			st.sections = parseDocSections(st.doc)

			// Extract struct fields.
			if astStruct, ok := ts.Type.(*ast.StructType); ok {
				for _, field := range astStruct.Fields.List {
					if len(field.Names) == 0 {
						continue
					}
					sf := structField{
						name:    field.Names[0].Name,
						comment: strings.TrimSpace(field.Comment.Text()),
					}
					if field.Tag != nil {
						sf.jsonTag = extractJSONTag(field.Tag.Value)
					}
					st.fields = append(st.fields, sf)
				}
			}

			structMap[ts.Name.Name] = st
		}
	}

	// Build result in order.
	var result []structType
	for _, name := range knownStructTypes {
		if st, ok := structMap[name]; ok {
			result = append(result, st)
		}
	}
	return result
}

func isKnownStruct(name string) bool {
	for _, n := range knownStructTypes {
		if n == name {
			return true
		}
	}
	return false
}

// parseDocSections splits a doc comment into the preamble (before any # heading)
// and sections delimited by lines starting with "# ".
func parseDocSections(doc string) []docSection {
	lines := strings.Split(doc, "\n")

	var sections []docSection
	var currentHeading string
	var currentBody []string
	inSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			// Save previous section if any.
			if inSection {
				sections = append(sections, docSection{
					heading: currentHeading,
					body:    strings.TrimSpace(strings.Join(currentBody, "\n")),
				})
			}
			currentHeading = strings.TrimPrefix(trimmed, "# ")
			currentBody = nil
			inSection = true
		} else if inSection {
			currentBody = append(currentBody, line)
		}
	}
	// Save last section.
	if inSection {
		sections = append(sections, docSection{
			heading: currentHeading,
			body:    strings.TrimSpace(strings.Join(currentBody, "\n")),
		})
	}

	return sections
}

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

func extractJSONTag(rawTag string) string {
	re := regexp.MustCompile(`json:"([^"]*)"`)
	m := re.FindStringSubmatch(rawTag)
	if len(m) < 2 {
		return ""
	}
	parts := strings.SplitN(m[1], ",", 2)
	return parts[0]
}

func cleanDoc(doc string) string {
	return strings.TrimSpace(doc)
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
	label string // display label, e.g. "system/init", "assistant(tool_use:Bash)"
	json  string // actual JSON output from MustJSON
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

// extractAssertPattern extracts label and JSON from a MustJSON(NewMessage*(...)) expression.
func extractAssertPattern(expr ast.Expr) (assertPattern, bool) {
	outer, ok := expr.(*ast.CallExpr)
	if !ok || len(outer.Args) != 1 {
		return assertPattern{}, false
	}
	outerIdent, ok := outer.Fun.(*ast.Ident)
	if !ok || outerIdent.Name != "MustJSON" {
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
	label, jsonStr := evalConstructor(innerIdent.Name, inner.Args)
	return assertPattern{label: label, json: jsonStr}, true
}

// evalConstructor maps a NewMessage* constructor name and its AST arguments
// to a display label and JSON string by actually calling the constructor.
func evalConstructor(funcName string, args []ast.Expr) (label, jsonStr string) {
	switch funcName {
	case "NewMessageSystemInit":
		return "system/init",
			ccprotocol.MustJSON(ccprotocol.NewMessageSystemInit())
	case "NewMessageSystemStatus":
		mode := extractStringLit(args[0])
		return "system/status(" + mode + ")",
			ccprotocol.MustJSON(ccprotocol.NewMessageSystemStatus(mode))
	case "NewMessageAssistantText":
		text := extractStringLit(args[0])
		return "assistant(text)",
			ccprotocol.MustJSON(ccprotocol.NewMessageAssistantText(text))
	case "NewMessageAssistantToolUse":
		name := extractStringLit(args[0])
		return "assistant(tool_use:" + name + ")",
			ccprotocol.MustJSON(ccprotocol.NewMessageAssistantToolUse(name))
	case "NewMessageAssistantThinking":
		thinking := extractStringLit(args[0])
		return "assistant(thinking)",
			ccprotocol.MustJSON(ccprotocol.NewMessageAssistantThinking(thinking))
	case "NewMessageUserToolResult":
		return "user(tool_result)",
			ccprotocol.MustJSON(ccprotocol.NewMessageUserToolResult())
	case "NewMessageUserToolResultError":
		return "user(tool_result:error)",
			ccprotocol.MustJSON(ccprotocol.NewMessageUserToolResultError())
	case "NewMessageResultSuccess":
		result := ""
		if len(args) > 0 {
			result = extractStringLit(args[0])
		}
		return "result/success",
			ccprotocol.MustJSON(ccprotocol.NewMessageResultSuccess(result))
	case "NewMessageResultSuccessIsError":
		return "result/success(is_error:true)",
			ccprotocol.MustJSON(ccprotocol.NewMessageResultSuccessIsError())
	case "NewMessageResultSuccessWithDenials":
		denials := extractPermissionDenials(args)
		return "result/success(permission_denials)",
			ccprotocol.MustJSON(ccprotocol.NewMessageResultSuccessWithDenials(denials...))
	case "NewMessageResultErrorDuringExecution":
		return "result/error_during_execution",
			ccprotocol.MustJSON(ccprotocol.NewMessageResultErrorDuringExecution())
	default:
		return funcName, "{}"
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

		for i, turn := range sc.turns {
			turnLabel := ""
			if len(sc.turns) > 1 {
				turnLabel = fmt.Sprintf(" %d", i+1)
			}
			inputJSON := `{"type":"user","message":{"role":"user","content":"..."}}`
			buf.WriteString("| ←" + turnLabel + " | [user](#message) | `" + inputJSON + "` |\n")

			for _, p := range turn {
				buf.WriteString("| → | [" + p.label + "](#message) | `" + p.json + "` |\n")
			}
		}

		buf.WriteString("\n")
	}
}

func writeSchemaSection(buf *strings.Builder, enums []enumType, structs []structType, msgFuncs []messageFunc) {
	buf.WriteString("## スキーマ\n\n")

	// Enum constants subsection.
	buf.WriteString("### enum 定数\n\n")
	for _, et := range enums {
		writeEnumType(buf, et)
	}

	// Type definitions subsection.
	buf.WriteString("### 型定義\n\n")
	for _, st := range structs {
		writeStructType(buf, st, msgFuncs)
	}
}

func writeEnumType(buf *strings.Builder, et enumType) {
	buf.WriteString("#### " + et.name + "\n\n")

	if et.doc != "" {
		buf.WriteString(et.doc + "\n\n")
	}

	if len(et.consts) > 0 {
		buf.WriteString("| 定数名 | 値 | 説明 |\n")
		buf.WriteString("|--------|-----|------|\n")
		for _, c := range et.consts {
			buf.WriteString(fmt.Sprintf("| `%s` | `\"%s\"` | %s |\n", c.name, c.value, c.comment))
		}
		buf.WriteString("\n")
	}
}

func writeStructType(buf *strings.Builder, st structType, msgFuncs []messageFunc) {
	buf.WriteString("#### " + st.name + "\n\n")

	// Write the preamble (doc text before the first # section).
	preamble := extractPreamble(st.doc)
	if preamble != "" {
		buf.WriteString(preamble + "\n\n")
	}

	// Fields table.
	if len(st.fields) > 0 {
		buf.WriteString("| フィールド | JSON キー | 説明 |\n")
		buf.WriteString("|-----------|-----------|------|\n")
		for _, f := range st.fields {
			buf.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n", f.name, f.jsonTag, f.comment))
		}
		buf.WriteString("\n")
	}

	// Render # type=... sections as subsections.
	for _, sec := range st.sections {
		buf.WriteString("##### " + sec.heading + "\n\n")
		buf.WriteString(godocToMarkdown(sec.body) + "\n\n")
	}

	// If this is Message, also write message constructor docs.
	if st.name == "Message" {
		for _, mf := range msgFuncs {
			buf.WriteString("##### " + mf.heading + "\n\n")
			buf.WriteString(godocToMarkdown(mf.body) + "\n\n")
		}
	}
}

// extractPreamble returns the doc text before the first "# " heading line.
func extractPreamble(doc string) string {
	lines := strings.Split(doc, "\n")
	var preambleLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			break
		}
		preambleLines = append(preambleLines, line)
	}
	return strings.TrimSpace(strings.Join(preambleLines, "\n"))
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

// skipFirstLine returns the doc text after the first line, trimmed.
func skipFirstLine(doc string) string {
	if idx := strings.Index(doc, "\n"); idx >= 0 {
		return strings.TrimSpace(doc[idx+1:])
	}
	return ""
}
