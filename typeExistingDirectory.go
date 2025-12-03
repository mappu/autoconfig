package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// ExistingDirectory allows browsing for an existing directory.
// The string value is the absolute path to the directory on disk.
type ExistingDirectory string

func (e ExistingDirectory) String() string {
	return string(e)
}

func (ExistingDirectory) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	hbox.AddWidget(rline.QWidget)

	browseBtn := qt.NewQPushButton2()
	setIcon(browseBtn.QAbstractButton, "folder-open", "Browse...", "Browse...")
	hbox.AddWidget(browseBtn.QWidget)

	browseBtn.OnClicked(func() {
		openDir := qt.QFileDialog_GetExistingDirectory3(browseBtn.QWidget, "Select a database directory...", rline.Text())
		if openDir != "" {
			rline.SetText(openDir)
		}
	})

	hboxWidget := qt.NewQWidget(area.ParentWidget())
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		rv.SetString(rline.Text())
	}
}
