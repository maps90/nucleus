package util

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestFromNullBool(t *testing.T) {
	test1 := sql.NullBool{Bool: true, Valid: true}
	result1 := FromNullBool(test1)
	assert.Equal(t, result1, true)

	test2 := sql.NullBool{Bool: false, Valid: false}
	result2 := FromNullBool(test2)
	assert.Equal(t, result2, false)
}

func TestFromNullInt(t *testing.T) {
	test1 := sql.NullInt64{Int64: 10, Valid: true}
	test2 := sql.NullInt64{Int64: 10, Valid: false}

	result1 := FromNullInt(test1)
	assert.Equal(t, result1, int64(10))

	result2 := FromNullInt(test2)
	assert.Equal(t, result2, int64(0))
}

func TestFromNullString(t *testing.T) {
	test1 := sql.NullString{String: "test", Valid: true}
	test2 := sql.NullString{String: "test", Valid: false}

	result1 := FromNullString(test1)
	assert.Equal(t, result1, "test")

	result2 := FromNullString(test2)
	assert.Equal(t, result2, "")
}

func TestFromNullFloat(t *testing.T) {
	test1 := sql.NullFloat64{Float64: 20.00, Valid: true}
	test2 := sql.NullFloat64{Float64: 20.00, Valid: false}

	result1 := FromNullFloat(test1)
	assert.Equal(t, result1, 20.00)

	result2 := FromNullFloat(test2)
	assert.Equal(t, result2, 0.00)

}

func TestFromNullTime(t *testing.T) {
	curTime := time.Now()
	var tm time.Time
	test1 := mysql.NullTime{Time: curTime, Valid: true}
	test2 := mysql.NullTime{Time: tm, Valid: false}

	result1 := FromNullTime(NullTime{test1})
	assert.Equal(t, result1, &curTime)

	result2 := FromNullTime(NullTime{test2})
	if result2 != nil {
		t.Error("unexpected results")
	}

}
