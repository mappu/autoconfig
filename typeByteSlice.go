package autoconfig

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

func handle_byte_slice(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	hbox := qt.NewQHBoxLayout2()

	display := qt.NewQLabel2()
	refreshDisplay := func() {
		content := rv.Bytes()
		if len(content) == 0 {
			display.SetText("Empty content")
			return
		}

		mimeType := http.DetectContentType(content)
		display.SetText(fmt.Sprintf("%s (%d bytes)", mimeType, len(content)))
	}
	display.SetSizePolicy2(qt.QSizePolicy__MinimumExpanding, qt.QSizePolicy__Minimum)
	refreshDisplay()
	hbox.AddWidget(display.QWidget)

	editBtn := qt.NewQToolButton2()

	menu := qt.NewQMenu(editBtn.QWidget)

	actionEditText := menu.AddActionWithText("Edit as text...")
	actionEditText.OnTriggered(func() {
		mlString := MultiLineString(rv.Bytes())
		OpenDialog(&mlString, editBtn.QWidget, label, func() {
			// Copy content from temp back into rv
			rv.SetBytes([]byte(mlString))
			refreshDisplay()
		})
	})

	actionLoadFromFile := menu.AddActionWithText("Import from file...")
	actionLoadFromFile.OnTriggered(func() {
		filePath := qt.QFileDialog_GetOpenFileNameWithParent(editBtn.QWidget)
		if filePath == "" {
			return // cancelled
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			qt.QMessageBox_Warning(editBtn.QWidget, "Error loading file content", err.Error())
			return
		}

		rv.SetBytes(content)
		refreshDisplay()
	})

	menu.AddSeparator()

	actionExport := menu.AddActionWithText("Export to file...")
	actionExport.OnTriggered(func() {
		filePath := qt.QFileDialog_GetSaveFileNameWithParent(editBtn.QWidget)
		if filePath == "" {
			return // cancelled
		}

		err := os.WriteFile(filePath, rv.Bytes(), 0644)
		if err != nil {
			qt.QMessageBox_Warning(editBtn.QWidget, "Error loading file content", err.Error())
			return
		}

		// Saved successfully
	})

	setIcon(editBtn.QAbstractButton, "document-edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	editBtn.SetMenu(menu)
	editBtn.SetPopupMode(qt.QToolButton__InstantPopup)
	hbox.AddWidget(editBtn.QWidget)

	addRowLayout(area, label, hbox.QLayout)

	return func() {
		// Edit function has already mutated the value
	}
}
