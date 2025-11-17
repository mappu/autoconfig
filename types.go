package autoconfig

import (
	"math"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

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

func handle_int(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rint := qt.NewQSpinBox2()

	// Range is split into upper+lower bounds
	var min, max int
	switch rv.Type().Bits() {
	case 8:
		min, max = math.MinInt8, math.MaxInt8
	case 16:
		min, max = math.MinInt16, math.MaxInt16
	case 32, 64:
		// QSpinBox is only capable of (signed) int32 maximum
		// TODO use a different widget
		min, max = math.MinInt32, math.MaxInt32
	}

	rint.SetMinimum(min)
	rint.SetMaximum(max)
	rint.SetValue(int(rv.Int())) // After setting bounds, otherwise it gets clamped

	area.AddRow3(label+`:`, rint.QWidget)
	return func() {
		rv.SetInt(int64(rint.Value()))
	}
}

func handle_uint(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rint := qt.NewQSpinBox2()
	// Range is entirely in nonnegative space

	rint.SetMinimum(0)
	switch rv.Type().Bits() {
	case 8:
		rint.SetMaximum(math.MaxUint8)
	case 16:
		rint.SetMaximum(math.MaxUint16)
	case 32, 64:
		// QSpinBox is only capable of (signed) int32 maximum
		// TODO use a different widget
		rint.SetMaximum(math.MaxInt32)
	}
	rint.SetValue(int(rv.Uint())) // After setting bounds, otherwise it gets clamped

	area.AddRow3(label+`:`, rint.QWidget)
	return func() {
		rv.SetUint(uint64(rint.Value()))
	}
}

func handle_float(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rfloat := qt.NewQDoubleSpinBox2()

	// By default, this is clamped to 100
	// Just allow ~unlimited, even for float32
	rfloat.SetMinimum(-math.MaxFloat64)
	rfloat.SetMaximum(math.MaxFloat64)
	rfloat.SetValue(rv.Float()) // After setting bounds, otherwise it gets clamped

	// This widget is also fixed to show two decimal places
	// May want to allow customization from a struct tag?

	area.AddRow3(label+`:`, rfloat.QWidget)
	return func() {
		rv.SetFloat(rfloat.Value())
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

func handle_pointer(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	statusField := qt.NewQLabel2()
	statusField.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Maximum)
	hbox.AddWidget(statusField.QWidget)

	refreshLabel := func() {
		statusField.SetText(formatValue(rv))
	}
	refreshLabel()

	configBtn := qt.NewQToolButton2()
	setIcon(configBtn.QAbstractButton, "edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	configBtn.OnClicked(func() {

		// Allocate our rv to be something if it's nothing
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))

			if defaulter, ok := rv.Interface().(InitDefaulter); ok {
				defaulter.InitDefaults()
			}
		}

		refreshLabel()

		// Going through .Interface() makes things non-addressible (Go cannot
		// assign through an interface).

		child := rv.Elem()

		openDialogFor(&child, configBtn.QWidget, label, func() {
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

func handle_stdlibTimeTime(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	rpicker := qt.NewQDateTimeEdit2()

	var ptrT *time.Time = (*time.Time)(rv.Addr().UnsafePointer())

	dt := qt.NewQDateTime()
	dt.SetSecsSinceEpoch(ptrT.Unix())

	rpicker.SetDateTime(dt)

	area.AddRow3(label+`:`, rpicker.QWidget)

	return func() {
		newTime := time.Unix(rpicker.DateTime().ToSecsSinceEpoch(), 0)
		*ptrT = newTime // assign
	}
}

func handle_slice(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	// HorizontalLayout
	// - QTreeWidget
	// - VerticalLayout
	//   - buttons x3
	//   - vspacer

	hbox := qt.NewQHBoxLayout2()

	itemList := qt.NewQTreeWidget2()
	itemList.SetHeaderHidden(true)
	itemList.SetUniformRowHeights(true)
	itemList.SetRootIsDecorated(false)
	itemList.SetSelectionMode(qt.QAbstractItemView__ContiguousSelection)
	hbox.AddWidget(itemList.QWidget)

	refreshListContent := func() {
		itemList.Clear()
		sliceItemsCt := rv.Len()
		for i := 0; i < sliceItemsCt; i++ {
			sliceElem := rv.Index(i)
			listItem := qt.NewQTreeWidgetItem2([]string{formatValue(&sliceElem)})
			itemList.AddTopLevelItem(listItem)
		}
	}
	refreshListContent()

	vbox := qt.NewQVBoxLayout2()
	// TODO if some buttons have icons created and some do not, the widths do
	// not match - should justify/align the widths(!)

	addButton := qt.NewQToolButton2()
	setIcon(addButton.QAbstractButton, "list-add", "+", "Add...")
	addButton.SetAutoRaise(true)
	addButton.OnClicked(func() {

		newElem := reflect.New(rv.Type().Elem() /* T */) // pointer-to-T, not a T
		if defaulter, ok := newElem.Interface().(InitDefaulter); ok {
			defaulter.InitDefaults()
		}

		openDialogFor(&newElem, addButton.QWidget, label, func() {

			// insert into slice
			maybeChangedRv := reflect.Append(*rv, newElem.Elem())
			rv.Set(maybeChangedRv)

			// refresh list
			refreshListContent()
		})
	})
	vbox.AddWidget(addButton.QWidget)

	editIndex := func(idx int) {
		curVal := rv.Index(idx)

		openDialogFor(&curVal, addButton.QWidget, label, func() {
			// we have directly mutated inside the slice already

			// refresh list
			refreshListContent()
		})
	}
	itemList.OnDoubleClicked(func(idx *qt.QModelIndex) {
		if idx == nil {
			return // doubleclick was not on an item
		}
		editIndex(idx.Row())
	})

	editButton := qt.NewQToolButton2()
	setIcon(editButton.QAbstractButton, "document-edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	editButton.SetAutoRaise(true)
	editButton.OnClicked(func() {
		curIdx := itemList.CurrentIndex()
		if curIdx == nil {
			return // nothing selected
		}
		editIndex(curIdx.Row())
	})

	vbox.AddWidget(editButton.QWidget)

	delButton := qt.NewQToolButton2()
	setIcon(delButton.QAbstractButton, "edit-delete-symbolic", "\u00d7" /* &times; */, "Remove")
	delButton.SetAutoRaise(true)

	delButton.OnClicked(func() {
		selectedItems := itemList.SelectedItems()

		// extract indexes
		var selectedIndexes []int = make([]int, 0, len(selectedItems))
		for _, itm := range selectedItems {
			selectedIndexes = append(selectedIndexes, itemList.IndexOfTopLevelItem(itm))
		}

		// reverse list, so indexes remain stable as we pop them
		sort.Slice(selectedIndexes, func(i, j int) bool { return i > j })

		// remove each item
		for _, removeIdx := range selectedIndexes {
			updated := rv.Slice(0, removeIdx)
			afterPart := rv.Slice(removeIdx+1, rv.Len())
			updated = reflect.AppendSlice(updated, afterPart)
			rv.Set(updated)
		}

		// re-render list
		refreshListContent()
	})

	vbox.AddWidget(delButton.QWidget)

	refreshButtonsEnabled := func() {
		selCt := len(itemList.SelectedItems())
		editButton.SetEnabled(selCt == 1)
		delButton.SetEnabled(selCt > 0)
	}
	refreshButtonsEnabled()
	itemList.OnSelectionChanged(func(super func(*qt.QItemSelection, *qt.QItemSelection), selected, deselected *qt.QItemSelection) {
		refreshButtonsEnabled()
		super(selected, deselected)
	})

	valign := qt.NewQSpacerItem4(0, 0, qt.QSizePolicy__Minimum, qt.QSizePolicy__MinimumExpanding)
	vbox.AddSpacerItem(valign)

	hbox.AddLayout(vbox.QLayout)

	area.AddRow4(label+`:`, hbox.QLayout)

	return func() {
	}
}
