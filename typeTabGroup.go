package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// If a TabGroup is the first member in a struct, the struct rendering will change
// to show all remaining fields as tabs.
//
// All the remaining struct members must be non-pointer types.
//
// Each struct member's 'ylabel' tag is used as the tab's label.
// Each struct member's 'yicon' tag, if present, is used as the tab's icon.
// The icon can either be a theme icon, or (with `:/` prefix) a Qt embedded resource.
type TabGroup struct{}

func handle_struct_as_TabGroup(area *qt.QFormLayout, rv *reflect.Value, _ reflect.StructTag, _ string) SaveFunc {

	obj := rv.Type()
	nf := obj.NumField()

	tabArea := qt.NewQTabWidget(area.ParentWidget())

	var allSavers []func()

	for i := 1; i < nf; i++ { // skip ourselves, we were element 0
		ff := obj.Field(i)  // Typeinfo only, not value
		valf := rv.Field(i) // Value

		// Handle icon

		useIcon := yicon_from_tag(ff.Tag)

		// Create tab frame

		frameWidget := qt.NewQWidget(area.ParentWidget())

		frame := qt.NewQFormLayout(frameWidget)

		// Don't pass in the struct's label here, we already showed it for the tab title
		saver := makeConfigAreaFor(&valf, frame, reflect.StructTag(""), "")

		if useIcon != nil {
			tabArea.AddTab2(frameWidget, useIcon, struct_field_label(ff))
		} else {
			tabArea.AddTab(frameWidget, struct_field_label(ff))
		}

		allSavers = append(allSavers, saver)
	}

	area.AddRowWithWidget(tabArea.QWidget)

	return func() {
		// Run all savers
		for _, saver := range allSavers {
			saver()
		}
	}
}
