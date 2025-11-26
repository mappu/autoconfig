package autoconfig

import (
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// If a OneOf is the first member in a struct, the struct rendering will change
// to allow selecting one of the remaining struct members.
// All the remaining struct members must be pointer types.
// Each struct member's 'ylabel' tag is used as the dropdown selection's label.
// When saving, only the selected struct member will be populated; all other
// values will be set to nil.
type OneOf string

func handle_struct_as_OneOf(area *qt.QFormLayout, rv *reflect.Value, _ reflect.StructTag, _ string) SaveFunc {

	obj := rv.Type()

	initialValue := rv.Field(0).String()
	var initialIndex int = 0

	picker := qt.NewQComboBox2()

	nf := obj.NumField()
	for i := 1; i < nf; i++ { // skip ourselves, we were element 0
		ff := obj.Field(i)

		picker.AddItem(struct_field_label(ff))
		if initialValue == ff.Name {
			initialIndex = i - 1
		}

		if iconTag, ok := ff.Tag.Lookup("yicon"); ok {
			// The icon might be a system theme icon ...
			if qt.QIcon_HasThemeIcon(iconTag) {
				picker.SetItemIcon(i-1, qt.QIcon_FromTheme(iconTag))

			} else if strings.HasPrefix(iconTag, `:/`) {
				// ... or it might be an embedded resource path
				picker.SetItemIcon(i-1, qt.NewQIcon4(iconTag))

			} else {
				// Shouldn't happen - probably the current PC has fewer
				// theme icons than the developer expected
				// No icon

			}
		}
	}

	picker.SetCurrentIndex(initialIndex)
	area.AddRowWithWidget(picker.QWidget)

	stack := qt.NewQStackedLayout2()

	var allSavers []func()

	for i := 1; i < nf; i++ { // skip ourselves, we were element 0
		ff := rv.Field(i)

		frameWidget := qt.NewQWidget2()

		frame := qt.NewQFormLayout2()
		frameWidget.SetLayout(frame.QLayout)

		if ff.Kind() != reflect.Pointer {
			// Weird, everything else in here should be a pointer
			panic("OneOf: expected all other struct members to be pointer types")
		}

		// If the value is nil, we have to new it, to have something to work with
		if ff.IsNil() {
			ff.Set(reflect.New(ff.Type().Elem()))
		}

		child := ff.Elem()

		saver := makeConfigAreaFor(&child, frame)

		stack.AddWidget(frameWidget)

		allSavers = append(allSavers, saver)
	}

	area.AddRowWithLayout(stack.QLayout)

	picker.OnCurrentIndexChanged(func(idx int) {
		stack.SetCurrentIndex(idx)
	})
	stack.SetCurrentIndex(initialIndex)

	return func() {

		// Commit current frame
		cidx := picker.CurrentIndex()
		allSavers[cidx]()

		// Save current selection into the picker value
		rv.Field(0).SetString(rv.Type().Field(cidx + 1).Name)

		for i := 1; i < nf; i++ {
			if i == (cidx + 1) {
				continue // keeping this one
			}

			// clearing this one
			ff := rv.Field(i)
			ff.SetZero()
		}
	}
}
