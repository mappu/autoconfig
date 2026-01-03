package autoconfig

import (
	qt "github.com/mappu/miqt/qt6"
)

// setIcon preferentially sets an icon for a button, using a label if the icon
// is not found.
func setIcon(btn *qt.QAbstractButton, iconThemeName, fallbackLabel, tooltip string) {
	if qt.QIcon_HasThemeIcon(iconThemeName) {
		btn.SetIcon(qt.QIcon_FromTheme(iconThemeName))
	} else {
		btn.SetText(fallbackLabel)
	}

	if fallbackLabel != tooltip {
		btn.SetToolTip(tooltip)
	}
}

// addRow adds the widget and label to the layout. It handles the case of a
// blank label.
func addRow(area *qt.QFormLayout, label string, widget *qt.QWidget) {
	if label == "" {
		area.AddRowWithWidget(widget) // No label
	} else {
		area.AddRow3(label+`:`, widget)
	}
}

// addRowLayout adds the widget and label to the layout. It handles the case of a
// blank label.
func addRowLayout(area *qt.QFormLayout, label string, layout *qt.QLayout) {
	if label == "" {
		area.AddRowWithLayout(layout) // No label
	} else {
		area.AddRow4(label+`:`, layout)
	}
}

// addRowBoxAndButtons adds the mainwidget and its buttons to the layout.
func addRowBoxAndButtons(area *qt.QFormLayout, label string, mainWidget *qt.QWidget, buttons ...*qt.QToolButton) {

	// HorizontalLayout
	// - QTreeWidget
	// - VerticalLayout
	//   - buttons x3
	//   - vspacer

	hbox := qt.NewQHBoxLayout2()
	hbox.AddWidget(mainWidget)

	vbox := qt.NewQVBoxLayout2()
	// TODO if some buttons have icons created and some do not, the widths do
	// not match - should justify/align the widths(!)

	for _, btn := range buttons {
		vbox.AddWidget(btn.QWidget)
	}

	valign := qt.NewQSpacerItem4(0, 0, qt.QSizePolicy__Minimum, qt.QSizePolicy__MinimumExpanding)
	vbox.AddSpacerItem(valign)

	hbox.AddLayout(vbox.QLayout)

	addRowLayout(area, label, hbox.QLayout)
}
