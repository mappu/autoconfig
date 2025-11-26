package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// Header shows a single-line header across the form.
// Use the `ylabel` tag to set the header's text.
type Header struct{}

func (Header) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rlabel := qt.NewQLabel3(label)
	area.AddRowWithWidget(rlabel.QWidget) // The widget spans both columns.
	return func() {}
}
