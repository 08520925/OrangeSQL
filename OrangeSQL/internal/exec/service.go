package exec

import (
	"OrangeSQL/internal/database"
)

// BuildResponse は database.ExecResult をAPIレスポンスに変換する。
func BuildResponse(result database.ExecResult, elapsedMs int64) Response {
	return Response{
		AffectedRows:  result.AffectedRows,
		ExecutionTime: elapsedMs,
	}
}
