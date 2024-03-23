package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

// FromNullString return a string if valid
// and empty string if not valid
func FromNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// FromNullFloat check if parameter is a valid float64
// if null return 0.00
func FromNullFloat(s sql.NullFloat64) float64 {
	if s.Valid {
		return s.Float64
	}
	return 0.00
}

//FromNullInt check if parameter is a valid int64
// if null return 0
func FromNullInt(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}

	return 0
}

//FromNullBool check if parameter is a valid boolean
// if null return false
func FromNullBool(i sql.NullBool) bool {
	if i.Valid {
		return i.Bool
	}

	return false
}

//FromNullTime check if parameter is a valid time
// return a pointer time
func FromNullTime(i sql.NullTime) *time.Time {
	if i.Valid {
		f := fmt.Sprintf("%v", i.Time)
		if f == "0001-01-01 00:00:00 +0000 UTC" {
			return nil
		}
		return &i.Time
	}

	return nil
}
