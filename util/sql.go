package util

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
)

// NullTime struct embbeded mysql.NullTime
type NullTime struct {
	mysql.NullTime
}

// FromNullString return a string if valid
// and empty string if not valid
func FromNullString(s sql.NullString) string {
	if s.Valid {
		return TrimWhiteSpace(s.String)
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
func FromNullTime(i NullTime) *time.Time {
	if i.Valid {
		f := fmt.Sprintf("%v", i.Time)
		if f == "0001-01-01 00:00:00 +0000 UTC" {
			return nil
		}
		return &i.Time
	}

	return nil
}

// NewNullString convert interface to sql.NullString
func NewNullString(s interface{}) (response sql.NullString) {
	casted, err := cast.ToStringE(s)
	if err != nil {
		response.Valid = false
		return response
	}

	response.String = casted
	response.Valid = true

	return
}

// NewNullInt convert interface to sql.NullInt64
func NewNullInt(s interface{}) (response sql.NullInt64) {
	casted, err := cast.ToInt64E(s)
	if err != nil {
		response.Valid = false
		return response
	}

	response.Int64 = casted
	response.Valid = true

	return
}

// NewNullBool convert interface to sql.NullBool
func NewNullBool(s interface{}) (response sql.NullBool) {
	casted, err := cast.ToBoolE(s)
	if err != nil {
		response.Valid = false
		return response
	}

	response.Bool = casted
	response.Valid = true

	return
}

// NewNullFloat convert interface to sql.NullFloat64
func NewNullFloat(s interface{}) (response sql.NullFloat64) {
	casted, err := cast.ToFloat64E(s)
	if err != nil {
		response.Valid = false
		return response
	}

	response.Float64 = casted
	response.Valid = true

	return
}
