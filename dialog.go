package autoconfig

import (
	qt "github.com/mappu/miqt/qt6"
)

// OpenDialog opens the struct for editing in a new modal dialog in the current
// global event loop.
// The dialog only has an "OK" button, you can't cancel your modifications to
// the supplied struct, the struct saver is always called.
func OpenDialog(ct ConfigurableStruct, parent *qt.QWidget, title string, onFinished func()) {

	dlg := qt.NewQDialog(parent)
	dlg.SetModal(true)
	// dlg.SetMinimumSize2(320, 240) // will grow
	dlg.SetWindowTitle(title)
	dlg.SetAttribute(qt.WA_DeleteOnClose)

	// QDialog
	// - VerticalLayout
	//   - FormLayout    <-- attach to config
	//   - QStandardButtonBar

	vbox := qt.NewQVBoxLayout(dlg.QWidget)
	vbox.SetContentsMargins(11, 11, 11, 11)
	vbox.SetSpacing(40)
	// dlg.SetLayout(vbox.QLayout)

	formArea := qt.NewQFormLayout2()
	formArea.SetContentsMargins(0, 0, 0, 0)
	formArea.SetSpacing(6)
	vbox.AddLayout(formArea.QLayout)
	applyer := MakeConfigArea(ct, formArea)

	buttons := qt.NewQDialogButtonBox(dlg.QWidget)
	buttons.SetStandardButtons(qt.QDialogButtonBox__Ok)
	buttons.OnAccepted(dlg.Accept)
	buttons.OnRejected(dlg.Reject)
	vbox.AddWidget(buttons.QWidget)

	dlg.OnFinished(func(status int) {
		// Save changes regardless of status
		applyer()
		onFinished()
	})

	dlg.Show()
}
