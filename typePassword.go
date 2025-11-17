package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// Password shows a single-line text area with character masking.
type Password string

func (Password) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQLineEdit2()
	rline.SetEchoMode(qt.QLineEdit__Password)
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}
