package sqlite

import (
	"database/sql"
	"time"
)

// toNullString returns sql.NullString. The string is considered valid if it's not empty.
func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// toNullBool returns sql.NullBool, that is always Valid
func toNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// toNullInt64 returns sql.NullInt64. The int is considered valid if it's not equal 0.
func toNullInt64(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

// toNullTime returns sql.NullTIme. The time is considered valid if it's not equal Zero.
func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}
