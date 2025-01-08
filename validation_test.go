package rqp

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	err := In("one", "two")("three")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "three: not in scope")

	err = In(1, 2)(3)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "3: not in scope")

	err = In(true)(false)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "false: not in scope")
}

func TestMinMax(t *testing.T) {
	err := Max(100)(101)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "101: not in scope")

	err = Max(100)(100)
	assert.NoError(t, err)

	err = Min(100)(100)
	assert.NoError(t, err)

	err = MinMax(10, 100)(9)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "9: not in scope")

	err = MinMax(10, 100)(101)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "101: not in scope")

	err = MinMax(10, 100)(50)
	assert.NoError(t, err)

	err = Multi(Min(10), Max(100))(50)
	assert.NoError(t, err)

	err = Multi(Min(10), Max(100))(101)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)

	err = MinMax(10, 100)("one")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "one: not in scope")

}

func TestMinMaxFloat(t *testing.T) {
	err := MaxFloat(100)(float32(101))
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "101: not in scope")

	err = MaxFloat(100)(float32(100))
	assert.NoError(t, err)

	err = MinFloat(100)(float32(100))
	assert.NoError(t, err)

	err = MinMaxFloat(10, 100)(float32(9))
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "9: not in scope")

	err = MinMaxFloat(10, 100)(float32(101))
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "101: not in scope")

	err = MinMaxFloat(10, 100)(float32(50))
	assert.NoError(t, err)

	err = Multi(MinFloat(10), MaxFloat(100))(float32(50))
	assert.NoError(t, err)

	err = Multi(MinFloat(10), MaxFloat(100))(float32(101))
	assert.Equal(t, errors.Cause(err), ErrNotInScope)

	err = MinMaxFloat(10, 100)("one")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "one: not in scope")
}

func TestMinMaxTime(t *testing.T) {
	current, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	assert.NoError(t, err)

	yesterday, err := time.Parse(time.RFC3339, "2018-12-31T00:00:00Z")
	assert.NoError(t, err)

	tomorrow, err := time.Parse(time.RFC3339, "2019-01-02T00:00:00Z")
	assert.NoError(t, err)

	err = MaxTime(current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02 00:00:00 +0000 UTC: not in scope")

	err = MaxTime(current)(current)
	assert.NoError(t, err)

	err = MinTime(current)(current)
	assert.NoError(t, err)

	err = MinMaxTime(current, tomorrow)(yesterday)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2018-12-31 00:00:00 +0000 UTC: not in scope")

	err = MinMaxTime(yesterday, current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02 00:00:00 +0000 UTC: not in scope")

	err = MinMaxTime(yesterday, tomorrow)(current)
	assert.NoError(t, err)

	err = Multi(MinTime(yesterday), MaxTime(tomorrow))(current)
	assert.NoError(t, err)

	err = Multi(MinTime(yesterday), MaxTime(current))(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)

	err = MinMaxTime(yesterday, tomorrow)("one")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "one: not in scope")
}

func TestMinMaxDate(t *testing.T) {
	current, err := civil.ParseDate("2019-01-01")
	assert.NoError(t, err)

	yesterday, err := civil.ParseDate("2018-12-31")
	assert.NoError(t, err)

	tomorrow, err := civil.ParseDate("2019-01-02")
	assert.NoError(t, err)

	err = MaxDate(current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02: not in scope")

	err = MaxDate(current)(current)
	assert.NoError(t, err)

	err = MinDate(current)(current)
	assert.NoError(t, err)

	err = MinMaxDate(current, tomorrow)(yesterday)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2018-12-31: not in scope")

	err = MinMaxDate(yesterday, current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02: not in scope")

	err = MinMaxDate(yesterday, tomorrow)(current)
	assert.NoError(t, err)

	err = Multi(MinDate(yesterday), MaxDate(tomorrow))(current)
	assert.NoError(t, err)

	err = Multi(MinDate(yesterday), MaxDate(current))(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)

	err = MinMaxDate(yesterday, tomorrow)("one")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "one: not in scope")
}

func TestMinMaxDateTime(t *testing.T) {
	current, err := civil.ParseDateTime("2019-01-01T00:00:00")
	assert.NoError(t, err)

	yesterday, err := civil.ParseDateTime("2018-12-31T00:00:00")
	assert.NoError(t, err)

	tomorrow, err := civil.ParseDateTime("2019-01-02T00:00:00")
	assert.NoError(t, err)

	err = MaxDateTime(current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02T00:00:00: not in scope")

	err = MaxDateTime(current)(current)
	assert.NoError(t, err)

	err = MinDateTime(current)(current)
	assert.NoError(t, err)

	err = MinMaxDateTime(current, tomorrow)(yesterday)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2018-12-31T00:00:00: not in scope")

	err = MinMaxDateTime(yesterday, current)(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "2019-01-02T00:00:00: not in scope")

	err = MinMaxDateTime(yesterday, tomorrow)(current)
	assert.NoError(t, err)

	err = Multi(MinDateTime(yesterday), MaxDateTime(tomorrow))(current)
	assert.NoError(t, err)

	err = Multi(MinDateTime(yesterday), MaxDateTime(current))(tomorrow)
	assert.Equal(t, errors.Cause(err), ErrNotInScope)

	err = MinMaxDateTime(yesterday, tomorrow)("one")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
	assert.EqualError(t, err, "one: not in scope")
}

func TestNotEmpty(t *testing.T) {
	// good case
	err := NotEmpty()("test")
	assert.NoError(t, err)
	// bad case
	err = NotEmpty()("")
	assert.Equal(t, errors.Cause(err), ErrNotInScope)
}
