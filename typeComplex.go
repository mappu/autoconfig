package autoconfig

import (
	"math"
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_complex(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	c := rv.Complex()

	// [input] + [input] i

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	// This widget is also fixed to show two decimal places
	// May want to allow customization from a struct tag?

	rep_float := qt.NewQDoubleSpinBox2()
	rep_float.SetMinimum(-math.MaxFloat64)
	rep_float.SetMaximum(math.MaxFloat64)
	rep_float.SetValue(real(c)) // After setting bounds, otherwise it gets clamped
	hbox.AddWidget(rep_float.QWidget)

	label1 := qt.NewQLabel3(`+`)
	hbox.AddWidget(label1.QWidget)

	imp_float := qt.NewQDoubleSpinBox2()
	imp_float.SetMinimum(-math.MaxFloat64)
	imp_float.SetMaximum(math.MaxFloat64)
	imp_float.SetValue(imag(c))
	imp_float.SetSuffix(" i")
	hbox.AddWidget(imp_float.QWidget)

	hboxWidget := qt.NewQWidget(area.ParentWidget())
	hboxWidget.SetLayout(hbox.QLayout)
	addRow(area, label, hboxWidget)

	return func() {
		rv.SetComplex(complex(rep_float.Value(), imp_float.Value()))
	}
}
