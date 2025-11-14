package autoconfig

import (
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
// Use the returned function to force all changes from the UI to be saved to
// the struct.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {

	obj := reflect.TypeOf(ct).Elem()

	var onApply []SaveFunc

	nf := obj.NumField()
	for i := 0; i < nf; i++ {
		i := i // go1.2xx

		ff := obj.Field(i)
		if !ff.IsExported() {
			continue
		}

		label := strings.ReplaceAll(ff.Name, `_`, ` `)   // Automatic name: field value with _ as spaces
		if useLabel, ok := ff.Tag.Lookup("ylabel"); ok { // Explicit name
			label = useLabel
		}

		widgetType := ff.Type.Name()

		// Maybe it is a struct pointer? If so, consider it an optional child dialog
		if ff.Type.Kind() == reflect.Pointer && ff.Type.Elem().Kind() == reflect.Struct {
			widgetType = "__childStruct"
		}

		handler, ok := registeredTypes[widgetType]
		if !ok {
			panic("makeConfigArea missing handling for type=" + widgetType)
		}

		fieldValue := reflect.ValueOf(ct).Elem().Field(i)

		singleFieldSaver := handler(area, &fieldValue, ff.Tag, label)

		onApply = append(onApply, func() {
			singleFieldSaver()
		})

	}

	// Save all
	return func() {
		for _, fn := range onApply {
			fn()
		}
	}
}
