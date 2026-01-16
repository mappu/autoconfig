package qspinbox

import (
	"math"
	"testing"
)

func TestSaturating(t *testing.T) {

	// Signed

	if got, want := AddSaturating(math.MaxInt64-100, 200), int64(math.MaxInt64); got != want { // Overflow
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := AddSaturating(math.MaxInt64-100, 50), int64(math.MaxInt64-50); got != want { // In bounds
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := AddSaturating(math.MinInt64+100, -200), int64(math.MinInt64); got != want { // Underflow
		t.Errorf("got %d, want %d", got, want)
	}

	// Unsigned

	if got, want := AddSaturatingUnsigned(math.MaxUint64-100, 200), uint64(math.MaxUint64); got != want { // Overflow
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := AddSaturatingUnsigned(math.MaxUint64-100, 50), uint64(math.MaxUint64-50); got != want { // In bounds
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := AddSaturatingUnsigned(100, -200), uint64(0); got != want { // Underflow
		t.Errorf("got %d, want %d", got, want)
	}
}
