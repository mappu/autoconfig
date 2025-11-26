package autoconfig

import (
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// EnumList allows choosing from a dropdown. The integer value is the 0-based index
// of available options.
// Available options should be set in the `yenum` struct tag, separated by ";;".
type EnumList int

func (EnumList) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	enumOpts, _ := tag.Lookup("yenum")

	rcombo := qt.NewQComboBox2()
	rcombo.AddItems(strings.Split(enumOpts, `;;`)) // Same separator as Qt filter (yfilter)
	rcombo.SetCurrentIndex(int(rv.Int()))

	area.AddRow3(label+`:`, rcombo.QWidget)

	return func() {
		rv.SetInt(int64(rcombo.CurrentIndex()))
	}
}
