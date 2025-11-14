package autoconfig

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

type ConfigurableStruct interface{}

type AddressPort struct {
	Address string
	Port    int
}

type ExistingFile string

type ExistingDirectory string

type Password string

type EnumList int

type saveHandler func()

type typeHandler func(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler

var (
	registeredTypes map[string]typeHandler
)

func handle_bool(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	rbtn := qt.NewQCheckBox3(label)
	rbtn.SetChecked(rv.Bool())
	area.AddRow3("", rbtn.QWidget)

	return func() {
		rv.SetBool(rbtn.IsChecked())
	}
}

func handle_string(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}

func handle_Password(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	rline := qt.NewQLineEdit2()
	rline.SetEchoMode(qt.QLineEdit__Password)
	rline.SetText(rv.String())
	area.AddRow3(label+`:`, rline.QWidget)
	return func() {
		rv.SetString(rline.Text())
	}
}

func handle_EnumList(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	enumOpts, _ := tag.Lookup("yenum")

	rcombo := qt.NewQComboBox2()
	rcombo.AddItems(strings.Split(enumOpts, `;;`)) // Same separator as Qt filter (yfilter)
	rcombo.SetCurrentIndex(int(rv.Int()))

	area.AddRow3(label+`:`, rcombo.QWidget)

	return func() {
		rv.SetInt(int64(rcombo.CurrentIndex()))
	}
}

func handle_ExistingFile(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	hbox.AddWidget(rline.QWidget)

	browseBtn := qt.NewQPushButton2()
	if qt.QIcon_HasThemeIcon("document-open") {
		browseBtn.SetIcon(qt.QIcon_FromTheme("document-open"))
		browseBtn.SetToolTip("Browse...")
	} else {
		browseBtn.SetText("Browse...")
	}

	hbox.AddWidget(browseBtn.QWidget)

	filter := "All files (*.*)"
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

func handle_ExistingDirectory(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	rline.SetText(rv.String())
	hbox.AddWidget(rline.QWidget)

	browseBtn := qt.NewQPushButton2()
	if qt.QIcon_HasThemeIcon("folder-open") {
		browseBtn.SetIcon(qt.QIcon_FromTheme("folder-open"))
		browseBtn.SetToolTip("Browse...")
	} else {
		browseBtn.SetText("Browse...")
	}
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

func handle_AddressPort(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {
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

func handle_ChildStructPtr(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) saveHandler {

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
	if qt.QIcon_HasThemeIcon("edit-symbolic") {
		configBtn.SetIcon(qt.QIcon_FromTheme("edit-symbolic"))
		configBtn.SetToolTip("Edit...")
	} else {
		configBtn.SetText("Edit...")
	}
	configBtn.OnClicked(func() {

		// Allocate our rv to be something if it's nothing
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
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
	if qt.QIcon_HasThemeIcon("edit-clear") {
		clearBtn.SetIcon(qt.QIcon_FromTheme("edit-clear"))
	} else {
		clearBtn.SetText("\u00d7") // &times; ×
	}
	clearBtn.SetToolTip("Clear")
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

func init() {
	registeredTypes = map[string]typeHandler{
		"bool":              handle_bool,
		"string":            handle_string,
		"Password":          handle_Password,
		"EnumList":          handle_EnumList,
		"ExistingFile":      handle_ExistingFile,
		"ExistingDirectory": handle_ExistingDirectory,
		"AddressPort":       handle_AddressPort,
		"__childStruct":     handle_ChildStructPtr,
	}
}
