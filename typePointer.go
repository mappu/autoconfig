package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_pointer(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	statusField := qt.NewQLabel2()
	statusField.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Maximum)
	hbox.AddWidget(statusField.QWidget)

	refreshLabel := func() {
		statusField.SetText(formatValue(rv))
	}
	refreshLabel()

	configBtn := qt.NewQToolButton2()
	setIcon(configBtn.QAbstractButton, "edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	configBtn.OnClicked(func() {

		// Allocate our rv to be something if it's nothing
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))

			if defaulter, ok := rv.Interface().(InitDefaulter); ok {
				defaulter.InitDefaults()
			}
		}

		refreshLabel()

		// Going through .Interface() makes things non-addressible (Go cannot
		// assign through an interface).

		child := rv.Elem()

		openDialogFor(&child, configBtn.QWidget, label, func() {
			// nothing to do
			refreshLabel()
		})
	})
	hbox.AddWidget(configBtn.QWidget)

	clearBtn := qt.NewQToolButton2()
	setIcon(clearBtn.QAbstractButton, "edit-clear", "\u00d7" /* &times; */, "Clear")
	clearBtn.OnClicked(func() {
		if !rv.IsNil() {
			rv.Set(reflect.Zero(rv.Type()))
		}
		refreshLabel()
	})
	hbox.AddWidget(clearBtn.QWidget)

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		// We have already mutated the *rv directly
	}
}
