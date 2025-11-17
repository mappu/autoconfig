package autoconfig

import (
	"reflect"
	"time"

	qt "github.com/mappu/miqt/qt6"
)

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
