package utils

// ---------------------------------------------------------------------------
// Any マッチャー — テストアサーション用のセンチネル値
// ---------------------------------------------------------------------------
// AssertOutput の完全一致比較で、動的な値（セッションID、所要時間など）を
// 型レベルでマッチさせるために使用する。各センチネルはフィールドの実際の型に
// 合わせた値を持ち、JSON ラウンドトリップ後に比較ロジックが認識する。

// AnyString は任意の JSON 値にマッチする。
// string 型フィールドにそのまま代入できる。
const AnyString = "<any>"

// AnyNumber は任意の JSON 数値にマッチする。
// float64 型フィールドに代入できる。値 -1 はセンチネルとして扱われる
// （duration, cost, turns 等は常に非負のため安全）。
var AnyNumber float64 = -1

// AnyStringSlice は任意の JSON 配列にマッチする。
// []string 型フィールドに代入できる。
var AnyStringSlice = []string{"<any>"}

// AnyMap は任意の JSON オブジェクトにマッチする。
// map[string]any 型フィールドに代入できる。
var AnyMap = map[string]any{"<any>": true}
