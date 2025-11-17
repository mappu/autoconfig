package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_bool(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rbtn := qt.NewQCheckBox3(label)
	rbtn.SetChecked(rv.Bool())
	area.AddRow3("", rbtn.QWidget)

	return func() {
		rv.SetBool(rbtn.IsChecked())
	}
}
