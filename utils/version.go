package utils

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// FieldMinVersion maps struct fields (StructName.FieldName) to the minimum CLI version
// that includes them. Used by test helpers to conditionally include version-specific
// fields in assertion patterns.
var FieldMinVersion = map[string]string{
	"SystemInitMessage.FastModeState": "2.1.38",
}

var (
	cliVersionOnce sync.Once
	cliVersionStr  string
)

// TestCLIVersion overrides the detected CLI version when non-empty.
// Used by gendoc to ensure all version-gated fields are included in documentation.
var TestCLIVersion string

// CLIVersion returns the installed Claude Code CLI version string (e.g. "2.1.38").
func CLIVersion() string {
	if TestCLIVersion != "" {
		return TestCLIVersion
	}
	cliVersionOnce.Do(func() {
		out, err := exec.Command("claude", "--version").Output()
		if err != nil {
			cliVersionStr = "0.0.0"
			return
		}
		// Output format: "2.1.38 (Claude Code)"
		cliVersionStr = strings.Fields(strings.TrimSpace(string(out)))[0]
	})
	return cliVersionStr
}

// CLIVersionAtLeast returns true if the installed CLI version is >= minVersion.
func CLIVersionAtLeast(minVersion string) bool {
	return compareVersions(CLIVersion(), minVersion) >= 0
}

// compareVersions compares two semver strings (major.minor.patch).
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) int {
	ap := parseVersion(a)
	bp := parseVersion(b)
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1
		}
		if ap[i] > bp[i] {
			return 1
		}
	}
	return 0
}

func parseVersion(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		result[i], _ = strconv.Atoi(parts[i])
	}
	return result
}
