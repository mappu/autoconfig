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
