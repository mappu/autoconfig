package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_fixed(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rlabel := qt.NewQLabel2()
	rlabel.SetText(formatValue(rv))
	area.AddRow3(label+`:`, rlabel.QWidget)
	return func() {}
}
