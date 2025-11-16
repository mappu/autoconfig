package autoconfig

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

type ConfigurableStruct interface{}

// InitDefaulter is a type that can reset itself to default values.
// It's used if autoconfig needs to initialize a child struct.
type InitDefaulter interface {
	InitDefaults()
}

type SaveFunc func()

// Autoconfiger is a custom-rendered type that can be interacted with
// automatically by the autoconfig package.
type Autoconfiger interface {
	Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc
}

func handle_bool(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rbtn := qt.NewQCheckBox3(label)
	rbtn.SetChecked(rv.Bool())
	area.AddRow3("", rbtn.QWidget)

	return func() {
		rv.SetBool(rbtn.IsChecked())
	}
}

func handle_string(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}

// MultiLineString shows a multi-line string area.
type MultiLineString string

func (MultiLineString) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQTextEdit2()
	rline.SetPlainText(rv.String())
	rline.SetAcceptRichText(false)
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.ToPlainText())
	}
}

// Password shows a single-line text area with character masking.
type Password string

func (Password) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rline := qt.NewQLineEdit2()
	rline.SetEchoMode(qt.QLineEdit__Password)
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}

// EnumList allows choosing from a dropdown. The integer value is the 0-based index
// of available options.
// Available options should be set in the `yenum` struct tag, separated by ";;".
type EnumList int

func (EnumList) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	enumOpts, _ := tag.Lookup("yenum")

	rcombo := qt.NewQComboBox2()
	rcombo.AddItems(strings.Split(enumOpts, `;;`)) // Same separator as Qt filter (yfilter)
	rcombo.SetCurrentIndex(int(rv.Int()))

	area.AddRow3(label+`:`, rcombo.QWidget)

	return func() {
		rv.SetInt(int64(rcombo.CurrentIndex()))
	}
}

// ExistingFile allows browsing for an existing file.
// The string value is the absolute path to the file on disk.
// If the `yfilter` struct tag is present, this allows constraining the file types using Qt syntax.
type ExistingFile string

func (ExistingFile) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
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

// ExistingDirectory allows browsing for an existing directory.
// The string value is the absolute path to the directory on disk.
type ExistingDirectory string

func (ExistingDirectory) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
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

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		rv.SetString(rline.Text())
	}
}

// AddressPort allows entering a text address and a numeric port.
// The port is limited to the 0-65535 range.
type AddressPort struct {
	Address string
	Port    int
}

func (AddressPort) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	addr := qt.NewQLineEdit2()
	addr.SetText(rv.Field(0).String()) // Address
	hbox.AddWidget(addr.QWidget)

	separator := qt.NewQLabel3(`:`)
	hbox.AddWidget(separator.QWidget)

	port := qt.NewQSpinBox2()
	port.SetMinimum(0)
	port.SetMaximum(65535)
	port.SetValue(int(rv.Field(1).Int())) // Port
	hbox.AddWidget(port.QWidget)

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		newVal := AddressPort{Address: addr.Text(), Port: port.Value()}
		rv.Set(reflect.ValueOf(newVal))
	}
}

func handle_ChildStructPtr(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	statusField := qt.NewQLabel2()
	statusField.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Maximum)
	hbox.AddWidget(statusField.QWidget)

	refreshLabel := func() {
		if rv.IsNil() {
			statusField.SetText("Not configured")
		} else if stringer, ok := rv.Interface().(fmt.Stringer); ok {
			statusField.SetText(stringer.String())
		} else {
			statusField.SetText("Configured")
		}
	}
	refreshLabel()

	configBtn := qt.NewQToolButton2()
	setIcon(configBtn.QAbstractButton, "edit-symbolic", "Edit...", "Edit...")
	configBtn.OnClicked(func() {

		// Allocate our rv to be something if it's nothing
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))

			if defaulter, ok := rv.Interface().(InitDefaulter); ok {
				defaulter.InitDefaults()
			}
		}

		refreshLabel()

		// Let OpenDialog mutate our new wipValue struct's fields directly
		OpenDialog(rv.Interface(), configBtn.QWidget, label, func() {
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

// Header shows a single-line header across the form.
// Use the `ylabel` tag to set the header's text.
type Header struct{}

func (Header) Autoconfig(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rlabel := qt.NewQLabel3(label)
	area.AddRowWithWidget(rlabel.QWidget) // The widget spans both columns.
	return func() {}
}
