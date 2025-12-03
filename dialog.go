package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// OpenDialog opens the struct for editing in a new modal dialog in the current
// global event loop.
// The dialog only has an "OK" button, you can't cancel your modifications to
// the supplied struct, the struct saver is always called.
func OpenDialog(ct ConfigurableStruct, parent *qt.QWidget, title string, onFinished func()) {
	rv := reflect.ValueOf(ct)
	openDialogFor(&rv, parent, title, onFinished)
}

func openDialogFor(rv *reflect.Value, parent *qt.QWidget, title string, onFinished func()) {

	dlg := qt.NewQDialog(parent)
	dlg.SetModal(true)
	dlg.SetWindowTitle(title)
	dlg.SetAttribute(qt.WA_DeleteOnClose)

	// QDialog
	// - VerticalLayout
	//   - FormLayout    <-- attach to config
	//   - QStandardButtonBar

	dlg.SetUpdatesEnabled(false) // Reduce flicker

	vbox := qt.NewQVBoxLayout(dlg.QWidget)
	vbox.SetContentsMargins(11, 11, 11, 11)
	vbox.SetSpacing(40)

	formArea := qt.NewQFormLayout2()
	formArea.SetContentsMargins(0, 0, 0, 0)
	formArea.SetSpacing(6)
	vbox.AddLayout(formArea.QLayout)
	applyer := makeConfigAreaFor(rv, formArea)

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

	dlg.SetUpdatesEnabled(true) // Reduce flicker

	dlg.Show()
}
