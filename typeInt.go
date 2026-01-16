package autoconfig

import (
	"math"
	"reflect"

	"github.com/mappu/autoconfig/qspinbox"
	qt "github.com/mappu/miqt/qt6"
)

func handle_int(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	switch rv.Type().Bits() {
	case 8:
		return handle_numeric_Int32SpinBox(area, rv, tag, label, math.MinInt8, math.MaxInt8)
	case 16:
		return handle_numeric_Int32SpinBox(area, rv, tag, label, math.MinInt16, math.MaxInt16)
	case 32:
		return handle_numeric_Int32SpinBox(area, rv, tag, label, math.MinInt32, math.MaxInt32)
	case 64:
		// Custom widget
		return handle_numeric_Int64SpinBox(area, rv, tag, label, math.MinInt64, math.MaxInt64)

	default:
		panic("Unknown bit width for integer type")
	}
}

func handle_uint(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	switch rv.Type().Bits() {
	case 8:
		return handle_numeric_Uint32SpinBox(area, rv, tag, label, math.MaxUint8)
	case 16:
		return handle_numeric_Uint32SpinBox(area, rv, tag, label, math.MaxUint16)

		// QSpinBox is only capable of (signed) int32 maximum
	case 32:
		// Custom widget
		return handle_numeric_Uint64SpinBox(area, rv, tag, label, math.MaxUint32)
	case 64:
		// Custom widget
		return handle_numeric_Uint64SpinBox(area, rv, tag, label, math.MaxUint64)

	default:
		panic("Unknown bit width for integer type")
	}
}

func handle_numeric_Int32SpinBox(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string, min int, max int) SaveFunc {
	rint := qt.NewQSpinBox2()
	rint.SetMinimum(min)
	rint.SetMaximum(max)
	rint.SetValue(int(rv.Int())) // After setting bounds, otherwise it gets clamped

	if prefix := tag.Get("yprefix"); len(prefix) > 0 {
		rint.SetPrefix(prefix)
	}
	if suffix := tag.Get("ysuffix"); len(suffix) > 0 {
		rint.SetSuffix(suffix)
	}

	addRow(area, label, rint.QWidget)
	return func() {
		rv.SetInt(int64(rint.Value()))
	}
}

func handle_numeric_Uint32SpinBox(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string, max int) SaveFunc {
	// WARNING: Only handles int32 bounds, not 0...uint32
	// Can be used for uint8/16, but, uint32/64 should both use the other specialized implementation

	rint := qt.NewQSpinBox2()
	rint.SetMinimum(0)
	rint.SetMaximum(max)
	rint.SetValue(int(rv.Uint())) // After setting bounds, otherwise it gets clamped

	if prefix := tag.Get("yprefix"); len(prefix) > 0 {
		rint.SetPrefix(prefix)
	}
	if suffix := tag.Get("ysuffix"); len(suffix) > 0 {
		rint.SetSuffix(suffix)
	}

	addRow(area, label, rint.QWidget)
	return func() {
		rv.SetUint(uint64(rint.Value()))
	}
}

func handle_numeric_Int64SpinBox(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string, min int64, max int64) SaveFunc {
	rint := qspinbox.NewQInt64SpinBox(nil)
	rint.SetMinimum(min)
	rint.SetMaximum(max)
	rint.SetValue(rv.Int()) // After setting bounds, otherwise it gets clamped

	if prefix := tag.Get("yprefix"); len(prefix) > 0 {
		rint.SetPrefix(prefix)
	}
	if suffix := tag.Get("ysuffix"); len(suffix) > 0 {
		rint.SetSuffix(suffix)
	}

	addRow(area, label, rint.QWidget)
	return func() {
		rv.SetInt(int64(rint.Value()))
	}
}

func handle_numeric_Uint64SpinBox(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string, max uint64) SaveFunc {
	rint := qspinbox.NewQUint64SpinBox(nil)
	rint.SetMinimum(0)
	rint.SetMaximum(max)
	rint.SetValue(rv.Uint()) // After setting bounds, otherwise it gets clamped

	if prefix := tag.Get("yprefix"); len(prefix) > 0 {
		rint.SetPrefix(prefix)
	}
	if suffix := tag.Get("ysuffix"); len(suffix) > 0 {
		rint.SetSuffix(suffix)
	}

	addRow(area, label, rint.QWidget)
	return func() {
		rv.SetUint(rint.Value())
	}
}
