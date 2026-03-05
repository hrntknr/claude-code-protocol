# fix-ci

## Overview

GA(GitHub Actions)の失敗を調査・修正し、PRを作成してCIが合格するまで繰り返すスキル。

## Workflow

1. **失敗の特定**: `gh run list --limit 5` で最新のCI実行を確認し、失敗しているジョブを特定する。
2. **ログ分析**: `gh api repos/{owner}/{repo}/actions/jobs/{job_id}/logs` で失敗ログを取得し、失敗テストとその原因を特定する。
3. **コード修正**: テストコードや関連コードを修正する。`go vet ./...` でビルド確認。
4. **ブランチ作成・コミット・プッシュ**:
   - `git checkout -b fix/<簡潔な説明>`
   - `git add <変更ファイル>` → `git commit -m "fix: <説明>"`
   - `git push -u origin <ブランチ>`
5. **PR作成**: `gh pr create --base main --title "fix: ..." --body "..."` でPRを作成。
6. **CI監視**: `gh pr checks <PR番号> --watch` でCIの完了を待つ。
   - 初回実行時はチェックが登録されるまで `sleep 15` してから `--watch` する。
7. **失敗時の再修正**: CIが失敗した場合:
   - 失敗ログを再取得・分析
   - コード修正 → コミット → プッシュ
   - 再度 `gh pr checks --watch` で確認
   - 全テストが合格するまで繰り返す
8. **完了報告**: 全チェック合格後、PR URLを報告。

## Notes

- ログ取得には `gh api repos/{owner}/{repo}/actions/jobs/{job_id}/logs` を使用する（`--log-failed` は出力が空になることがある）。
- `gh pr checks --watch` の初回はチェック未登録エラーになるため `sleep 15` を挟む。
- 複数バージョンでテストされている場合、全バージョンで合格する修正が必要。
- バージョン間で動作が異なる場合は `utils.CLIVersionAtLeast()` で分岐するか、共通部分のみアサートする。
- 非決定的な動作がある場合は、決定的な部分のみをテストする方針にする。
