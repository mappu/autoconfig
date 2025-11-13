package autoconfig

import (
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// assignInterfaceStructField allows you to modify a struct's field by ordinal.
func assignInterfaceStructField(target ConfigurableStruct, fieldId int, cb func(*reflect.Value)) {
	// Values contained in an interface are not addressable
	// Copy the struct value to a temporary variable, set the field
	// in the temporary variable and copy the temporary variable
	// back to the interface.

	// v is the interface{}
	v := reflect.ValueOf(&target).Elem()

	// Allocate a temporary variable with type of the struct.
	//    v.Elem() is the value contained in the interface.
	tmp := reflect.New(v.Elem().Type()).Elem()

	// Copy the struct value contained in interface to
	// the temporary variable.
	tmp.Set(v.Elem())

	// Set the field.
	// setText := rline.Text()
	// tmp.Elem().Field(i) //.SetString(setText)
	field := tmp.Elem().Field(fieldId)
	cb(&field)

	// Set the interface to the modified struct value.
	v.Set(tmp)
}

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
func MakeConfigArea(ct ConfigurableStruct, area *qt.QFormLayout) func() ConfigurableStruct {

	obj := reflect.TypeOf(ct).Elem()

	return makeConfigAreaFor(obj, area)
}

// MakeConfigArea makes a config area by pushing elements into a QFormLayout.
func makeConfigAreaFor(obj reflect.Type, area *qt.QFormLayout) func() ConfigurableStruct {

	makeAssigner := func(fieldId int, cb func(*reflect.Value)) func(ConfigurableStruct) {
		return func(target ConfigurableStruct) {
			assignInterfaceStructField(target, fieldId, cb)
		}
	}

	var onApply []func(ConfigurableStruct)

	nf := obj.NumField()
	for i := 0; i < nf; i++ {
		i := i // go1.2xx

		ff := obj.Field(i)
		if !ff.IsExported() {
			continue
		}

		label := strings.ReplaceAll(ff.Name, `_`, ` `)   // Automatic name: field value with _ as spaces
		if useLabel, ok := ff.Tag.Lookup("ylabel"); ok { // Explicit name
			label = useLabel
		}

		widgetType := ff.Type.Name()

		// Maybe it is a struct pointer? If so, consider it an optional child dialog
		if ff.Type.Kind() == reflect.Pointer && ff.Type.Elem().Kind() == reflect.Struct {
			widgetType = "__childStruct"
		}

		switch widgetType {
		case "bool":
			rbtn := qt.NewQCheckBox3(label)
			area.AddRow3("", rbtn.QWidget)
			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				rv.SetBool(rbtn.IsChecked())
			}))

		case "string", "Password":
			rline := qt.NewQLineEdit2()
			if widgetType == "Password" {
				rline.SetEchoMode(qt.QLineEdit__Password)

			}
			if useInit, ok := ff.Tag.Lookup("yinit"); ok {
				rline.SetText(useInit)
			}
			area.AddRow3(label+`:`, rline.QWidget)
			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				rv.SetString(rline.Text())
			}))

		case "EnumList":
			enumOpts, _ := ff.Tag.Lookup("yenum")

			rcombo := qt.NewQComboBox2()
			rcombo.AddItems(strings.Split(enumOpts, `;;`)) // Same separator as Qt filter (yfilter)

			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				rv.SetInt(int64(rcombo.CurrentIndex()))
			}))

			area.AddRow3(label+`:`, rcombo.QWidget)

		case "ExistingFile", "ExistingDirectory":
			hbox := qt.NewQHBoxLayout2()
			hbox.SetContentsMargins(0, 0, 0, 0)

			rline := qt.NewQLineEdit2()
			hbox.AddWidget(rline.QWidget)

			if useInit, ok := ff.Tag.Lookup("yinit"); ok {
				rline.SetText(useInit)
			}

			browseBtn := qt.NewQPushButton3("Browse...")
			hbox.AddWidget(browseBtn.QWidget)

			filter := "All files (*.*)"
			if useFilter, ok := ff.Tag.Lookup("yfilter"); ok {
				filter = useFilter
			}

			if widgetType == "ExistingDirectory" {
				browseBtn.OnClicked(func() {
					openDir := qt.QFileDialog_GetExistingDirectory3(browseBtn.QWidget, "Select a database directory...", rline.Text())
					if openDir != "" {
						rline.SetText(openDir)
					}
				})

			} else if widgetType == "ExistingFile" {
				browseBtn.OnClicked(func() {
					startDir := filepath.Dir(rline.Text())

					openPath := qt.QFileDialog_GetOpenFileName4(browseBtn.QWidget, "Select a database file...", startDir, filter)
					if openPath != "" {
						rline.SetText(openPath)
					}
				})
			}

			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				rv.SetString(rline.Text())
			}))

			hboxWidget := qt.NewQWidget2()
			hboxWidget.SetLayout(hbox.QLayout)
			area.AddRow3(label+`:`, hboxWidget)

		case "AddressPort":
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

			if defaultVal, ok := ff.Tag.Lookup(`yport`); ok {
				defaultValInt, err := strconv.Atoi(defaultVal)
				if err != nil {
					panic(err)
				}
				port.SetValue(defaultValInt)
			}

			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				newVal := AddressPort{Address: addr.Text(), Port: port.Value()}
				rv.Set(reflect.ValueOf(newVal))
			}))

			hboxWidget := qt.NewQWidget2()
			hboxWidget.SetLayout(hbox.QLayout)
			area.AddRow3(label+`:`, hboxWidget)

		case "__childStruct":
			hbox := qt.NewQHBoxLayout2()
			hbox.SetContentsMargins(0, 0, 0, 0)

			statusField := qt.NewQLabel2()
			statusField.SetText("Not configured")
			statusField.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Maximum)
			hbox.AddWidget(statusField.QWidget)

			configBtn := qt.NewQToolButton2()
			configBtn.SetText("Edit...")
			configBtn.OnClicked(func() {

				// Allocate a temporary variable with type of the struct.
				//    v.Elem() is the value contained in the interface.

				openDialogFor(ff.Type.Elem(), configBtn.QWidget, label, func(childThing ConfigurableStruct) {
					// ...
				})
			})
			hbox.AddWidget(configBtn.QWidget)

			clearBtn := qt.NewQToolButton2()
			clearBtn.SetText("\u00d7") // &times; ×
			hbox.AddWidget(clearBtn.QWidget)

			onApply = append(onApply, makeAssigner(i, func(rv *reflect.Value) {
				// ...
			}))

			hboxWidget := qt.NewQWidget2()
			hboxWidget.SetLayout(hbox.QLayout)
			area.AddRow3(label+`:`, hboxWidget)

		default:
			panic("makeConfigArea missing handling for type=" + widgetType)
		}

	}

	getter := func() ConfigurableStruct {

		// Get a zero-valued version of the struct to start with
		var ct ConfigurableStruct = reflect.New(obj).Interface().(ConfigurableStruct)

		for _, fn := range onApply {
			fn(ct)
		}
		return ct
	}
	return getter
}
