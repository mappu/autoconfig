package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

const defaultLabel = "Configure"

type ConfigurableStruct interface{}

// Resetter is a type that can reset itself to default values.
// It's used if autoconfig needs to initialize a child struct.
type Resetter interface {
	Reset()
}

type SaveFunc func()

// Renderer is a custom-rendered type that can be interacted with
// automatically by the autoconfig package.
type Renderer interface {
	Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc
}

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
// Use the returned function to force all changes from the UI to be saved to
// the struct.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) SaveFunc {

	rv := reflect.ValueOf(ct)
	return makeConfigAreaFor(&rv, area, reflect.StructTag(""), defaultLabel)
}

func makeConfigAreaFor(rv *reflect.Value, area *qt.QFormLayout, tag reflect.StructTag, label string) SaveFunc {

	// If this layout is already placed inside a widget, reduce layout reflow flicker
	// This seems to only affect Windows, not Linux
	if pwdg := area.ParentWidget(); pwdg != nil && pwdg.UpdatesEnabled() {
		pwdg.SetUpdatesEnabled(false)
		defer pwdg.SetUpdatesEnabled(true)
	}

	return handle_any(area, rv, tag, label)
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

	if rv.Type().Kind() == reflect.Pointer {
		// Handle before any other cases (Renderer)
		// If this is a pointer type, we always want it to go the 'Optional' style
		return handle_pointer(area, rv, tag, label)

	} else if renderer, ok := rv.Interface().(Renderer); ok {
		// The Renderer interface implemented with a Value receiver and we have a value
		return renderer.Render(area, rv, tag, label)

	} else if renderer, ok := rv.Addr().Interface().(Renderer); ok {
		// The Renderer interface implemented with a Pointer receiver and we have a value
		return renderer.Render(area, rv, tag, label)

	} else if rv.Type().String() == "time.Time" {
		return handle_stdlibTimeTime(area, rv, tag, label) // Handle this case earlier, otherwise, it would match Struct

	} else {
		switch rv.Type().Kind() {
		case reflect.Func, reflect.UnsafePointer, reflect.Chan:
			// No way we can configure these types
			return handle_fixed(area, rv, tag, label)

		case reflect.Bool:
			return handle_bool(area, rv, tag, label)

		case reflect.String:
			return handle_string(area, rv, tag, label)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return handle_int(area, rv, tag, label)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return handle_uint(area, rv, tag, label)

		case reflect.Float32, reflect.Float64:
			return handle_float(area, rv, tag, label)

		case reflect.Struct:
			// Struct by non-pointer
			// Integrate it directly
			return handle_struct(area, rv, tag, label)

		case reflect.Slice:
			return handle_slice(area, rv, tag, label)

		case reflect.Pointer:
			return handle_pointer(area, rv, tag, label)

		case reflect.Interface:
			// If it's an interface (error, io.Reader, io.Writer, ...) then skip it
			return handle_fixed(area, rv, tag, label)

		case reflect.Complex64,
			reflect.Complex128,
			reflect.Array,
			reflect.Map:
			// TODO
			// These are probably representable but not yet implemented
			return handle_fixed(area, rv, tag, label)

		default:
			// The above enum should have covered every constant Kind available
			// in the stdlib reflect package
			// If there's something new in here, either data is corrupt or
			// a future version of Go has added something fundamentally new
			panic("makeConfigArea missing handling for type=" + rv.Type().String())
		}
	}

}
