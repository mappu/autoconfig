package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// MultiLineString shows a multi-line string area.
type MultiLineString string

func (MultiLineString) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQTextEdit2()
	rline.SetPlainText(rv.String())
	rline.SetAcceptRichText(false)
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.ToPlainText())
	}
}
