package autoconfig

import (
	"reflect"

	qt "github.com/mappu/miqt/qt6"
)

// EnumString allows choosing from a dropdown.
// First, prefill the available options via the SetEnumStringOptions function,
// and then pass your global keyname in the `yenum` struct tag.
type EnumString string

var enumStringOpts map[string][]string

// SetEnumStringOptions configures the list of allowed options for the given key
// when used with the autoconfig.EnumString type.
// This is package-global and not threadsafe.
// To unregister a key, set the 'options' to nil.
func SetEnumStringOptions(key string, options []string) {
	if enumStringOpts == nil {
		enumStringOpts = make(map[string][]string)
	}

	if options == nil {
		delete(enumStringOpts, key)
	} else {
		enumStringOpts[key] = options
	}
}

func (EnumString) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	enumKey, _ := tag.Lookup("yenum")

	opts, ok := enumStringOpts[enumKey]
	if !ok {
		// programmer error
		panic("EnumString: key '" + enumKey + "' not registered in SetEnumListOptions")
	}

	currentIndex := 0
	{
		currentString := rv.String()
		for i, opt := range opts {
			if opt == currentString {
				currentIndex = i
				break
			}
		}
	}

	rcombo := qt.NewQComboBox2()
	rcombo.AddItems(opts)
	rcombo.SetCurrentIndex(currentIndex)

	addRow(area, label, rcombo.QWidget)

	return func() {
		rv.SetString(opts[rcombo.CurrentIndex()])
	}
}
