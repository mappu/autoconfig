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

	if rv.Type().Kind() == reflect.Func {
		// No way we can configure a function

	} else if autoconfiger, ok := rv.Interface().(Autoconfiger); ok {
		handler = autoconfiger.Autoconfig

	} else if rv.Type().String() == "bool" {
		handler = handle_bool

	} else if rv.Type().String() == "string" {
		handler = handle_string

	} else if rv.Type().String() == "time.Time" {
		handler = handle_stdlibTimeTime // Handle this case earlier, otherwise, it would match Struct

	} else if rv.Type().Kind() == reflect.Struct {
		// Struct by non-pointer
		// Integrate it directly
		handler = handle_struct

	} else if rv.Type().Kind() == reflect.Slice {
		handler = handle_slice

	} else if rv.Type().Kind() == reflect.Pointer {
		handler = handle_pointer

	} else if rv.Type().Kind() == reflect.Interface {
		// If it's an interface (error, io.Reader, io.Writer, ...) then skip it

	} else {
		// A real unsupported type
		panic("makeConfigArea missing handling for type=" + rv.Type().String())
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
