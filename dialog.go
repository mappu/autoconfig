package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// OpenDialog opens the struct for editing in a new modal dialog in the current
// global event loop.
// The onFinished callback receives nil on cancel.
func OpenDialog(ct ConfigurableStruct, parent *qt.QWidget, title string, onFinished func(ConfigurableStruct)) {

	obj := reflect.TypeOf(ct).Elem()

	openDialogFor(obj, parent, title, onFinished)
}

// OpenDialog opens the struct for editing in a new modal dialog in the current
// global event loop.
// The onFinished callback receives nil on cancel.
func openDialogFor(obj reflect.Type, parent *qt.QWidget, title string, onFinished func(ConfigurableStruct)) {

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
	applyer := makeConfigAreaFor(obj, formArea)

	buttons := qt.NewQDialogButtonBox(dlg.QWidget)
	buttons.SetStandardButtons(qt.QDialogButtonBox__Ok | qt.QDialogButtonBox__Cancel)
	buttons.OnAccepted(dlg.Accept)
	buttons.OnRejected(dlg.Reject)
	vbox.AddWidget(buttons.QWidget)

	dlg.OnFinished(func(status int) {
		if status != int(qt.QDialog__Accepted) {
			onFinished(nil)
			return
		}

		updatedStruct := applyer()
		onFinished(updatedStruct)
	})

	dlg.Show()
}
