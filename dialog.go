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
	openDialogFor(&rv, parent, reflect.StructTag(""), title, onFinished)
}

func openDialogFor(rv *reflect.Value, parent *qt.QWidget, tag reflect.StructTag, title string, onFinished func()) {

	dlg := qt.NewQDialog(parent)
	dlg.SetModal(true)
	dlg.SetWindowTitle(title)
	dlg.SetAttribute(qt.WA_DeleteOnClose)
	dlg.SetUpdatesEnabled(false) // Reduce flicker

	// QDialog
	// - VerticalLayout
	//   - QScrollArea
	//     - QWidget (viewport)
	//       - FormLayout    <-- attach to config
	//   - QStandardButtonBar

	formArea := qt.NewQFormLayout2()
	formArea.SetContentsMargins(0, 0, 0, 0)
	formArea.SetSpacing(6)
	formArea.SetSizeConstraint(qt.QLayout__SetMinAndMaxSize)
	// Pass through a blank label. The main label is in the dialog header instead.
	applyer := makeConfigAreaFor(rv, formArea, tag, "")

	viewport := qt.NewQWidget(dlg.QWidget)
	viewport.SetLayout(formArea.QLayout)

	scrollArea := qt.NewQScrollArea2()
	scrollArea.SetFrameShape(qt.QFrame__NoFrame)
	scrollArea.SetWidgetResizable(true)
	szp := scrollArea.VerticalScrollBar().SizePolicy()
	szp.SetRetainSizeWhenHidden(true)
	scrollArea.VerticalScrollBar().SetSizePolicy(*szp)
	scrollArea.SetWidget(viewport)

	vbox := qt.NewQVBoxLayout(dlg.QWidget)
	vbox.SetContentsMargins(11, 11, 11, 11)
	vbox.SetSpacing(40)
	vbox.AddWidget(scrollArea.QWidget)

	buttons := qt.NewQDialogButtonBox(dlg.QWidget)
	buttons.SetStandardButtons(qt.QDialogButtonBox__Ok)
	buttons.OnAccepted(dlg.Accept)
	buttons.OnRejected(dlg.Reject)
	vbox.AddWidget(buttons.QWidget)

	dlg.SetLayout(vbox.QLayout)

	dlg.OnFinished(func(status int) {
		// Save changes regardless of status
		applyer()
		onFinished()
	})

	dlg.SetUpdatesEnabled(true) // Reduce flicker

	dlg.Show()

	// Ensure the dialog always fits an additional vertical scrollbar, without causing
	// a horizontal scrollbar to also appear
	// FIXME This causes horizontal flicker when opening the dialog, but it can't
	// be done before .Show()
	// TODO replace with real value from current style metrics
	const ESTIMATE_VSCROLLBAR_WIDTH = 32

	dlg.SetMinimumWidth(dlg.Width() + ESTIMATE_VSCROLLBAR_WIDTH)
}
