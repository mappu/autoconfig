package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
// Use the returned function to force all changes from the UI to be saved to
// the struct.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {

	rv := reflect.ValueOf(ct)
	return makeConfigAreaFor(&rv, area)
}

func makeConfigAreaFor(rv *reflect.Value, area *qt.QFormLayout) SaveFunc {
	return handle_any(area, rv, reflect.StructTag(""), "Configure")
}

func handle_any(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	if !rv.CanAddr() {
		// Sometimes we'll be supplied with something not addressible, but, points
		// to something that is addressible
		if rv.Kind() == reflect.Pointer && rv.Elem().CanAddr() {
			// Use that instead
			child := rv.Elem()
			return handle_any(area, &child, tag, label)
		}

		panic("Supplied value is not addressable, cannot be mutated?")
	}

	type typeHandler func(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc

	var handler typeHandler = nil

	if autoconfiger, ok := rv.Interface().(Autoconfiger); ok {
		handler = autoconfiger.Autoconfig

	} else if rv.Type().String() == "time.Time" {
		handler = handle_stdlibTimeTime // Handle this case earlier, otherwise, it would match Struct

	} else {
		switch rv.Type().Kind() {
		case reflect.Func, reflect.UnsafePointer, reflect.Chan:
			// No way we can configure these types
			handler = handle_fixed

		case reflect.Bool:
			handler = handle_bool

		case reflect.String:
			handler = handle_string

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			handler = handle_int

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			handler = handle_uint

		case reflect.Float32, reflect.Float64:
			handler = handle_float

		case reflect.Struct:
			// Struct by non-pointer
			// Integrate it directly
			handler = handle_struct

		case reflect.Slice:
			handler = handle_slice

		case reflect.Pointer:
			handler = handle_pointer

		case reflect.Interface:
			// If it's an interface (error, io.Reader, io.Writer, ...) then skip it
			handler = handle_fixed

		case reflect.Complex64,
			reflect.Complex128,
			reflect.Array,
			reflect.Map:
			// TODO
			// These are probably representable but not yet implemented
			handler = handle_fixed

		default:
			// The above enum should have covered every constant Kind available
			// in the stdlib reflect package
			// If there's something new in here, either data is corrupt or
			// a future version of Go has added something fundamentally new
			panic("makeConfigArea missing handling for type=" + rv.Type().String())
		}
	}

	return handler(area, rv, tag, label)

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

		label := formatLabel(ff.Name)                    // Automatic name: field value with _ as spaces
		if useLabel, ok := ff.Tag.Lookup("ylabel"); ok { // Explicit name
			label = useLabel
		}

		fieldValue := rv.Field(i)

		singleFieldSaver := handle_any(area, &fieldValue, ff.Tag, label)

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
