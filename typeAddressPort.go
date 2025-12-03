package autoconfig

import (
	"fmt"
	"reflect"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// AddressPort allows entering a text address and a numeric port.
// The port is limited to the 0-65535 range.
type AddressPort struct {
	Address string
	Port    int
}

func (AddressPort) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
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

	hboxWidget := qt.NewQWidget(area.ParentWidget())
	hboxWidget.SetLayout(hbox.QLayout)
	area.AddRow3(label+`:`, hboxWidget)

	return func() {
		newVal := AddressPort{Address: addr.Text(), Port: port.Value()}
		rv.Set(reflect.ValueOf(newVal))
	}
}

func (a AddressPort) String() string {
	if a == (AddressPort{}) {
		return "<Not set>"
	}

	if strings.Contains(a.Address, `:`) {
		return fmt.Sprintf("[%s]:%d", a.Address, a.Port)
	}

	return fmt.Sprintf("%s:%d", a.Address, a.Port)
}
