package utils

import (
	"strings"
)

// Pattern holds a JSON assertion pattern with paths to ignore during comparison.
type Pattern struct {
	json        string
	ignorePaths map[string]bool
}

// NewPattern creates a Pattern from a JSON string with optional ignore paths.
func NewPattern(json string, ignorePaths ...string) Pattern {
	m := make(map[string]bool, len(ignorePaths))
	for _, p := range ignorePaths {
		m[p] = true
	}
	return Pattern{json: json, ignorePaths: m}
}

// Ignore returns a new Pattern with additional paths to ignore during comparison.
func (p Pattern) Ignore(paths ...string) Pattern {
	m := make(map[string]bool, len(p.ignorePaths)+len(paths))
	for k, v := range p.ignorePaths {
		m[k] = v
	}
	for _, path := range paths {
		m[path] = true
	}
	return Pattern{json: p.json, ignorePaths: m}
}

// Assert returns a new Pattern with paths removed from the ignore set,
// enabling comparison for those paths.
func (p Pattern) Assert(paths ...string) Pattern {
	m := make(map[string]bool, len(p.ignorePaths))
	for k, v := range p.ignorePaths {
		m[k] = v
	}
	for _, path := range paths {
		delete(m, path)
	}
	return Pattern{json: p.json, ignorePaths: m}
}

// String returns the JSON string (implements fmt.Stringer).
func (p Pattern) String() string {
	return p.json
}

// isIgnored checks if the given path should be ignored during comparison.
// Supports wildcard * matching for array element paths.
func isIgnored(ignorePaths map[string]bool, path string) bool {
	if ignorePaths[path] {
		return true
	}
	for pattern := range ignorePaths {
		if strings.Contains(pattern, "*") && matchWildcard(pattern, path) {
			return true
		}
	}
	return false
}

// matchWildcard matches a dot-separated path pattern with * wildcards
// against a concrete path. Each * matches exactly one path segment.
func matchWildcard(pattern, path string) bool {
	pp := strings.Split(pattern, ".")
	tp := strings.Split(path, ".")
	if len(pp) != len(tp) {
		return false
	}
	for i, seg := range pp {
		if seg == "*" {
			continue
		}
		if seg != tp[i] {
			return false
		}
	}
	return true
}
