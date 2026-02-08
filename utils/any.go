package utils

// ---------------------------------------------------------------------------
// Any matchers — sentinel values for test assertions
// ---------------------------------------------------------------------------
// Used in AssertOutput's exact matching to match dynamic values (session ID,
// duration, etc.) at the type level. Each sentinel holds a value matching the
// field's actual type, recognized by the comparison logic after JSON round-trip.

// AnyString matches any JSON value.
// Can be assigned directly to string fields.
const AnyString = "<any>"

// AnyNumber matches any JSON number.
// Can be assigned to float64 fields. The value -1 is used as a sentinel
// (safe because duration, cost, turns, etc. are always non-negative).
var AnyNumber float64 = -1

// AnyStringSlice matches any JSON array.
// Can be assigned to []string fields.
var AnyStringSlice = []string{"<any>"}

// AnyMap matches any JSON object.
// Can be assigned to map[string]any fields.
var AnyMap = map[string]any{"<any>": true}
