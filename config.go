package autoconfig

import (
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// assignInterfaceStructField allows you to modify a struct's field by ordinal.
func assignInterfaceStructField(target ConfigurableStruct, fieldId int, cb func(*reflect.Value)) {
	// Values contained in an interface are not addressable
	// Copy the struct value to a temporary variable, set the field
	// in the temporary variable and copy the temporary variable
	// back to the interface.

	// v is the interface{}
	v := reflect.ValueOf(&target).Elem()

	// Allocate a temporary variable with type of the struct.
	//    v.Elem() is the value contained in the interface.
	tmp := reflect.New(v.Elem().Type()).Elem()

	// Copy the struct value contained in interface to
	// the temporary variable.
	tmp.Set(v.Elem())

	// Set the field.
	// setText := rline.Text()
	// tmp.Elem().Field(i) //.SetString(setText)
	field := tmp.Elem().Field(fieldId)
	cb(&field)

	// Set the interface to the modified struct value.
	v.Set(tmp)
}

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) func() ConfigurableStruct {

	obj := reflect.TypeOf(ct).Elem()

	return makeConfigAreaFor(obj, area)
}

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
func makeConfigAreaFor(obj reflect.Type, area *qt.QFormLayout) func() ConfigurableStruct {

	makeAssigner := func(fieldId int, cb func(*reflect.Value)) func(ConfigurableStruct) {
		return func(target ConfigurableStruct) {
			assignInterfaceStructField(target, fieldId, cb)
		}
	}

	var onApply []func(ConfigurableStruct)

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

		singleFieldSaver := handler(area, ff.Type, ff.Tag, label)

		onApply = append(onApply, makeAssigner(i, singleFieldSaver))

	}

	getter := func() ConfigurableStruct {

		// Get a zero-valued version of the struct to start with
		var ct ConfigurableStruct = reflect.New(obj).Interface().(ConfigurableStruct)

		for _, fn := range onApply {
			fn(ct)
		}
		return ct
	}
	return getter
}
