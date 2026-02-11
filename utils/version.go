package utils

import (
	"encoding/json"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// FieldMinVersion maps "StructName.json_key" to the minimum CLI version
// that includes the field. Used by MustJSONVersioned to omit version-gated
// fields when running against older CLIs.
var FieldMinVersion = map[string]string{
	"SystemInitMessage.fast_mode_state": "2.1.38",
}

// MustJSONVersioned marshals v to JSON and removes fields whose CLI version
// requirement (per FieldMinVersion) is not met by the installed CLI.
// Keys in FieldMinVersion use the format "StructName.json_key".
func MustJSONVersioned(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic("MustJSONVersioned: " + err.Error())
	}

	structName := reflect.TypeOf(v).Name()
	cur := CLIVersion()

	var omitKeys []string
	for key, minVer := range FieldMinVersion {
		sn, jsonKey, ok := strings.Cut(key, ".")
		if !ok || sn != structName {
			continue
		}
		if !CLIVersionAtLeast(cur, minVer) {
			omitKeys = append(omitKeys, jsonKey)
		}
	}

	if len(omitKeys) == 0 {
		return string(b)
	}

	var raw map[string]json.RawMessage
	json.Unmarshal(b, &raw)
	for _, k := range omitKeys {
		delete(raw, k)
	}
	b, _ = json.Marshal(raw)
	return string(b)
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

// CLIVersionAtLeast returns true if cur is >= minVersion.
func CLIVersionAtLeast(cur, minVersion string) bool {
	return compareVersions(cur, minVersion) >= 0
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
