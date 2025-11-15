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

		type typeHandler func(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc

		var handler typeHandler = nil

		if autoconfiger, ok := reflect.ValueOf(ct).Elem().Field(i).Interface().(Autoconfiger); ok {
			handler = autoconfiger.Autoconfig

		} else if ff.Type.Kind() == reflect.Pointer && ff.Type.Elem().Kind() == reflect.Struct {
			// Maybe it is a struct pointer? If so, consider it an optional child dialog
			handler = handle_ChildStructPtr

		} else {
			// Hardcoded implementations for builtin types
			switch ff.Type.Name() {
			case "bool":
				handler = handle_bool
			case "string":
				handler = handle_string
			default:
				panic("makeConfigArea missing handling for type=" + ff.Type.Name())
			}
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
