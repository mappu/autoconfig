package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
// Use the returned function to force all changes from the UI to be saved to
// the struct.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {

	typ := reflect.TypeOf(ct)

	if typ.Kind() == reflect.Struct {
		// struct by value
		// This doesn't work, makeConfigAreaForStruct() will immediately call .Elem()!
		panic("struct by value, expected by pointer")

	} else if typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct {
		// struct by pointer (still works)
		return makeConfigAreaForStruct(ct, area)

	} else if typ.Kind() == reflect.Pointer {
		// Recurse
		return makeConfigAreaForPointer(ct, area)

	} else {
		// singular/non-struct/the struct is deeper than a single pointer level
		panic(reflect.TypeOf(ct).String())
	}
}

func makeConfigAreaForPointer(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {
	obj := reflect.ValueOf(ct).Elem()

	return handle_ChildStructPtr(area, &obj, reflect.StructTag(""), formatLabel(obj.Type().String()))
}

func makeConfigAreaForStruct(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {
	obj := reflect.ValueOf(ct).Elem()

	return handle_struct(area, &obj, reflect.StructTag(""), formatLabel(obj.Type().String()))
}

func handle_struct(area *qt.QFormLayout, rv *reflect.Value, _ reflect.StructTag, _ string) SaveFunc {

	// ignore tag and label

	obj := rv.Type()

	var onApply []SaveFunc

	nf := obj.NumField()
	for i := 0; i < nf; i++ {
		i := i // go1.2xx

		ff := obj.Field(i)

		// Don't show private fields
		if !ff.IsExported() {
			continue
		}

		if ff.Type.Kind() == reflect.Func {
			continue // No way we can configure a function
		}

		label := formatLabel(ff.Name)                    // Automatic name: field value with _ as spaces
		if useLabel, ok := ff.Tag.Lookup("ylabel"); ok { // Explicit name
			label = useLabel
		}

		type typeHandler func(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc

		var handler typeHandler = nil

		if autoconfiger, ok := rv.Field(i).Interface().(Autoconfiger); ok {
			handler = autoconfiger.Autoconfig

		} else if ff.Type.Kind() == reflect.Pointer && ff.Type.Elem().Kind() == reflect.Struct {
			// Maybe it is a struct pointer? If so, consider it an optional child dialog
			handler = handle_ChildStructPtr

		} else if ff.Type.Kind() == reflect.Struct {
			// Struct by non-pointer
			// Integrate it directly
			handler = handle_struct

		} else if ff.Type.Kind() == reflect.Slice {
			handler = handle_slice

		} else if ff.Type.Kind() == reflect.Pointer {
			handler = handle_ChildStructPtr

		} else {
			// Hardcoded implementations for builtin types
			switch ff.Type.String() {
			case "bool":
				handler = handle_bool
			case "string":
				handler = handle_string
			case "time.Time":
				handler = handle_stdlibTimeTime
			default:
				// If it's an interface (error, io.Reader, io.Writer, ...) then
				// skip it
				if ff.Type.Kind() == reflect.Interface {
					continue
				}

				// A real unsupported type
				panic("makeConfigArea missing handling for type=" + ff.Type.String())
			}
		}

		fieldValue := rv.Field(i)

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
