package autoconfig

import (
	"math"
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_float(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rfloat := qt.NewQDoubleSpinBox2()

	// By default, this is clamped to 100
	// Just allow ~unlimited, even for float32
	rfloat.SetMinimum(-math.MaxFloat64)
	rfloat.SetMaximum(math.MaxFloat64)
	rfloat.SetValue(rv.Float()) // After setting bounds, otherwise it gets clamped

	// This widget is also fixed to show two decimal places
	// May want to allow customization from a struct tag?

	area.AddRow3(label+`:`, rfloat.QWidget)
	return func() {
		rv.SetFloat(rfloat.Value())
	}
}
