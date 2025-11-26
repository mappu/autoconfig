package autoconfig

import (
	"path/filepath"
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// ExistingFile allows browsing for an existing file.
// The string value is the absolute path to the file on disk.
// If the `yfilter` struct tag is present, this allows constraining the file types using Qt syntax.
type ExistingFile string

func (*ExistingFile) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	hbox.AddWidget(rline.QWidget)

	browseBtn := qt.NewQPushButton2()
	setIcon(browseBtn.QAbstractButton, "document-open", "Browse...", "Browse...")

	hbox.AddWidget(browseBtn.QWidget)

	filter := "All files (*)"
	if useFilter, ok := tag.Lookup("yfilter"); ok {
		filter = useFilter
	}

	browseBtn.OnClicked(func() {
		startDir := filepath.Dir(rline.Text())

		openPath := qt.QFileDialog_GetOpenFileName4(browseBtn.QWidget, "Select a database file...", startDir, filter)
		if openPath != "" {
			rline.SetText(openPath)
		}
	})

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		rv.SetString(rline.Text())
	}
}
