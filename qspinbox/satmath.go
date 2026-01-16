package qspinbox

import (
	"math"
)

// AddSaturating returns start+delta, clamping overflow/underflow to int64 bounds.
func AddSaturating(start int64, delta int) int64 {
	steps_i64 := int64(delta)

	if steps_i64 < 0 && start+steps_i64 > start {
		return math.MinInt64

	} else if steps_i64 > 0 && start+steps_i64 < start {
		return math.MaxInt64

	} else {
		return start + steps_i64

	}
}

// AddSaturatingUnsigned returns start+delta, clamping overflow/underflow to uint64 bounds.
func AddSaturatingUnsigned(start uint64, delta int) uint64 {

	if delta > 0 {
		if start+uint64(delta) < start {
			// The steps are positive, but adding them resulted in us being below the original value = overflow
			return math.MaxUint64
		}
		return start + uint64(delta)

	} else if delta < 0 {
		if start-uint64(-delta) > start {
			// The steps are negative, but subtracting(abs) resulted in us being higher than the original value = underflow
			return 0
		}
		return start - uint64(-delta)

	} else {
		return start
	}
}
