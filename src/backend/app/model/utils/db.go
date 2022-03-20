package utils

// 需要pgsql执行完sql后返回的字段
// http://www.postgresql.org/docs/current/static/sql-insert.html
// http://www.postgresql.org/docs/current/static/sql-update.html
// http://www.postgresql.org/docs/current/static/sql-delete.html
//
func SqlReturning(s string) string {
	if len(s) == 0 {
		return ""
	}
	return "RETURNING " + s
}

func SqlColumnsComplementTable(table string, columns ...string) []string {
	if len(columns) == 0 {
		return []string{}
	}
	result := make([]string, len(columns))
	for i := range columns {
		result[i] = table + "." + columns[i]
	}
	return result
}
