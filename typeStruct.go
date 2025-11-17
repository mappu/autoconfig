package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

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
