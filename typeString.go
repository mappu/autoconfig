package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_string(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}
