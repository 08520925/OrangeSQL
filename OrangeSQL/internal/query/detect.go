package query

import (
	"regexp"
	"strings"

	"OrangeSQL/internal/database"
)

// 単一テーブル SELECT を判定する正規表現。
// FROM <tableName> を抽出し、JOIN / UNION / サブクエリを除外する。
var fromPattern = regexp.MustCompile(`(?i)\bFROM\s+([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)?)`)

// 編集不可にする SQL キーワード。
var nonEditablePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bJOIN\b`),
	regexp.MustCompile(`(?i)\bUNION\b`),
	regexp.MustCompile(`(?i)\bGROUP\s+BY\b`),
	regexp.MustCompile(`(?i)\bHAVING\b`),
	regexp.MustCompile(`(?i)\bDISTINCT\b`),
}

// DetectEditableTable は SQL を解析して編集可能なテーブル名と PK カラムを返す。
// 編集不可の場合は空文字と nil を返す。
func DetectEditableTable(sql string, resultColumns []string, db database.Database) (tableName string, pkColumns []string) {
	trimmed := strings.TrimSpace(sql)

	// SELECT 文のみ対象
	if !strings.HasPrefix(strings.ToUpper(trimmed), "SELECT") {
		return "", nil
	}

	// サブクエリ検出: FROM の前に ( があるか、複数の FROM があるか
	if strings.Count(strings.ToUpper(trimmed), "FROM") > 1 {
		return "", nil
	}

	// 編集不可パターンチェック
	for _, p := range nonEditablePatterns {
		if p.MatchString(trimmed) {
			return "", nil
		}
	}

	// FROM <tableName> を抽出
	match := fromPattern.FindStringSubmatch(trimmed)
	if match == nil {
		return "", nil
	}
	table := match[1]

	// テーブルのカラム情報を取得して PK を特定
	cols, err := db.Columns(table)
	if err != nil {
		return "", nil
	}

	var pks []string
	for _, col := range cols {
		if col.PK {
			pks = append(pks, col.Name)
		}
	}

	// PK がなければ編集不可
	if len(pks) == 0 {
		return "", nil
	}

	// PK カラムが結果に含まれているか確認
	resultColSet := make(map[string]bool, len(resultColumns))
	for _, c := range resultColumns {
		resultColSet[strings.ToLower(c)] = true
	}
	for _, pk := range pks {
		if !resultColSet[strings.ToLower(pk)] {
			return "", nil
		}
	}

	return table, pks
}
