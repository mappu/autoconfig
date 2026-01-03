package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func struct_field_label(ff reflect.StructField) string {
	if useLabel, ok := ff.Tag.Lookup("ylabel"); ok { // Explicit name
		return useLabel
	}

	// Automatic name: field value with _ as spaces
	return formatLabel(ff.Name)
}

func handle_struct(area *qt.QFormLayout, rv *reflect.Value, self_tag reflect.StructTag, self_label string) SaveFunc {

	// ignore tag and label

	obj := rv.Type()

	var onApply []SaveFunc

	nf := obj.NumField()
	for i := 0; i < nf; i++ {
		ff := obj.Field(i)

		// Hooks for special struct types
		if i == 0 {
			switch {
			case ff.Type == reflect.TypeOf(OneOf("")):
				return handle_struct_as_OneOf(area, rv, self_tag, self_label)
			case ff.Type == reflect.TypeOf(TabGroup{}):
				return handle_struct_as_TabGroup(area, rv, self_tag, self_label)
			}
		}

		// Don't show private fields
		if !ff.IsExported() {
			continue
		}

		fieldValue := rv.Field(i)

		singleFieldSaver := handle_any(area, &fieldValue, ff.Tag, struct_field_label(ff))

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
