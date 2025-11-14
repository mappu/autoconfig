package autoconfig

import (
	"log"
	"path/filepath"
	"reflect"
	"strconv"
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

type saveHandler func(rv *reflect.Value)

type typeHandler func(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler

var (
	registeredTypes map[string]typeHandler
)

func handle_bool(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	rbtn := qt.NewQCheckBox3(label)
	area.AddRow3("", rbtn.QWidget)

	return func(rv *reflect.Value) {
		rv.SetBool(rbtn.IsChecked())
	}
}

func handle_string(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	rline := qt.NewQLineEdit2()
	if useInit, ok := tag.Lookup("yinit"); ok {
		rline.SetText(useInit)
	}
	area.AddRow3(label+`:`, rline.QWidget)
	return func(rv *reflect.Value) {
		rv.SetString(rline.Text())
	}
}

func handle_Password(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	rline := qt.NewQLineEdit2()
	rline.SetEchoMode(qt.QLineEdit__Password)
	if useInit, ok := tag.Lookup("yinit"); ok {
		rline.SetText(useInit)
	}
	area.AddRow3(label+`:`, rline.QWidget)
	return func(rv *reflect.Value) {
		rv.SetString(rline.Text())
	}
}

func handle_EnumList(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	enumOpts, _ := tag.Lookup("yenum")

	rcombo := qt.NewQComboBox2()
	rcombo.AddItems(strings.Split(enumOpts, `;;`)) // Same separator as Qt filter (yfilter)

	area.AddRow3(label+`:`, rcombo.QWidget)

	return func(rv *reflect.Value) {
		rv.SetInt(int64(rcombo.CurrentIndex()))
	}
}

func handle_ExistingFile(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	hbox.AddWidget(rline.QWidget)

	if useInit, ok := tag.Lookup("yinit"); ok {
		rline.SetText(useInit)
	}

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

	return func(rv *reflect.Value) {
		rv.SetString(rline.Text())
	}
}

func handle_ExistingDirectory(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rline := qt.NewQLineEdit2()
	hbox.AddWidget(rline.QWidget)

	if useInit, ok := tag.Lookup("yinit"); ok {
		rline.SetText(useInit)
	}

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

	return func(rv *reflect.Value) {
		rv.SetString(rline.Text())
	}
}

func handle_AddressPort(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {
	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	addr := qt.NewQLineEdit2()
	hbox.AddWidget(addr.QWidget)

	separator := qt.NewQLabel3(`:`)
	hbox.AddWidget(separator.QWidget)

	port := qt.NewQSpinBox2()
	hbox.AddWidget(port.QWidget)
	port.SetMinimum(0)
	port.SetMaximum(65535)

	if defaultVal, ok := tag.Lookup(`yport`); ok {
		defaultValInt, err := strconv.Atoi(defaultVal)
		if err != nil {
			panic(err)
		}
		port.SetValue(defaultValInt)
	}

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func(rv *reflect.Value) {
		newVal := AddressPort{Address: addr.Text(), Port: port.Value()}
		rv.Set(reflect.ValueOf(newVal))
	}
}

func handle_ChildStructPtr(area *qt.QFormLayout, typ reflect.Type, tag reflect.StructTag, label string) saveHandler {

	// Allocate a temporary variable with type of the struct.
	wipValue := reflect.New(typ.Elem()) // struct itself, not pointer
	isAllocated := false
	// But then new'ing it has given us a pointer again

	log.Printf("wipValue type = %q", wipValue.Type().String())
	log.Printf("wipValue.Elem() type = %q", wipValue.Elem().Type().String())
	wipValue.Elem().Field(0).SetBool(true)

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	statusField := qt.NewQLabel2()
	statusField.SetText("Not configured")
	statusField.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Maximum)
	hbox.AddWidget(statusField.QWidget)

	configBtn := qt.NewQToolButton2()
	if qt.QIcon_HasThemeIcon("edit-symbolic") {
		configBtn.SetIcon(qt.QIcon_FromTheme("edit-symbolic"))
		configBtn.SetToolTip("Edit...")
	} else {
		configBtn.SetText("Edit...")
	}
	configBtn.OnClicked(func() {

		openDialogFor(typ.Elem(), configBtn.QWidget, label, func(childThing ConfigurableStruct) {
			if childThing == nil {
				// Cancelled, do not modify our current wipValue

			} else {
				// childThing is interface
				iface := reflect.ValueOf(childThing)

				wipValue.Elem().Set(iface.Elem())

				statusField.SetText("Configured")
				isAllocated = true
			}
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
		statusField.SetText("Not configured")
		isAllocated = false
	})
	hbox.AddWidget(clearBtn.QWidget)

	hboxWidget := qt.NewQWidget2()
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func(rv *reflect.Value) {
		if isAllocated {
			rv.Set(wipValue)
		}
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
