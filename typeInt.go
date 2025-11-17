package autoconfig

import (
	"math"
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_int(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rint := qt.NewQSpinBox2()

	// Range is split into upper+lower bounds
	var min, max int
	switch rv.Type().Bits() {
	case 8:
		min, max = math.MinInt8, math.MaxInt8
	case 16:
		min, max = math.MinInt16, math.MaxInt16
	case 32, 64:
		// QSpinBox is only capable of (signed) int32 maximum
		// TODO use a different widget
		min, max = math.MinInt32, math.MaxInt32
	}

	rint.SetMinimum(min)
	rint.SetMaximum(max)
	rint.SetValue(int(rv.Int())) // After setting bounds, otherwise it gets clamped

	area.AddRow3(label+`:`, rint.QWidget)
	return func() {
		rv.SetInt(int64(rint.Value()))
	}
}

func handle_uint(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rint := qt.NewQSpinBox2()
	// Range is entirely in nonnegative space

	rint.SetMinimum(0)
	switch rv.Type().Bits() {
	case 8:
		rint.SetMaximum(math.MaxUint8)
	case 16:
		rint.SetMaximum(math.MaxUint16)
	case 32, 64:
		// QSpinBox is only capable of (signed) int32 maximum
		// TODO use a different widget
		rint.SetMaximum(math.MaxInt32)
	}
	rint.SetValue(int(rv.Uint())) // After setting bounds, otherwise it gets clamped

	area.AddRow3(label+`:`, rint.QWidget)
	return func() {
		rv.SetUint(uint64(rint.Value()))
	}
}
