package autoconfig

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mappu/autoconfig/qspinbox"
	qt "github.com/mappu/miqt/qt6"
)

// Factor is an int64 where the effective value is multiplied by some factor.
//
// Use with the `yfactor` tag, e.g.
//
//	yfactor:"1;;B;;1024;;KiB;;1048576;;MiB;;1073741824;;GiB"
//
// Rules:
//   - There must be an even number of fields, separated by double-semicolon (;;).
//   - The first value of each pair must be an integer.
//   - There should usually be a first entry set to `1` so that a sensible value
//     is available for any possible integer input.
//   - They should be sorted smallest to largest, but this is not enforced.
//
// The Factor renderer will automatically pick the right-most type that divides
// the input value without any remainder.
//
// Negative values are allowed.
//
// It is possible to use this to construct an int64 overflow (large input with
// large factor)... Just don't do that.
type Factor int64

type factor struct {
	Divisor int64
	Label   string
}

func (Factor) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	// Parse all factors out

	yfactor := tag.Get("yfactor")
	if len(yfactor) == 0 {
		// Factor without yfactor tag is just an int64
		return handle_int(area, rv, tag, label)
	}

	parts := strings.Split(yfactor, `;;`)
	if len(parts)%2 != 0 {
		panic("autoconfig.Factor expects yfactor to have an even number of properties") // Programmer error
	}

	factors := make([]factor, 0, len(parts)/2)

	for i := 0; i < len(parts); i += 2 {
		parse, err := strconv.ParseInt(parts[i], 10, 64)
		if err != nil {
			panic(err) // Programmer error
		}

		factors = append(factors, factor{parse, parts[i+1]})
	}

	return handle_factor_with(area, rv, tag, label, factors)
}

// handle_factor_with is the common helper for Factor-type inputs.
func handle_factor_with(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string, factors []factor) SaveFunc {

	// Determine current factor for input value

	rawValue := rv.Int()
	initialFactorIdx := 0
	for i := 0; i < len(factors); i += 1 {
		reverseIdx := len(factors) - i - 1
		if rawValue%factors[reverseIdx].Divisor == 0 {
			initialFactorIdx = reverseIdx // OK, match
			break
		}
	}

	// Construct

	hbox := qt.NewQHBoxLayout2()
	hbox.SetContentsMargins(0, 0, 0, 0)

	rint := qspinbox.NewQInt64SpinBox(nil)
	rint.SetMinimum(math.MinInt64)
	rint.SetMaximum(math.MaxInt64) // Send bounds first, otherwise SetValue() gets clamped
	rint.SetValue(rv.Int() / factors[initialFactorIdx].Divisor)
	hbox.AddWidget(rint.QWidget)

	opts := qt.NewQComboBox2()
	for _, fac := range factors {
		opts.AddItem(fac.Label)
	}
	opts.SetCurrentIndex(initialFactorIdx)
	hbox.AddWidget(opts.QWidget)

	hboxWidget := qt.NewQWidget(area.ParentWidget())
	hboxWidget.SetLayout(hbox.QLayout)
	addRow(area, label, hboxWidget)

	return func() {
		rawValue := int64(rint.Value()) * factors[opts.CurrentIndex()].Divisor
		rv.SetInt(rawValue)
	}

}

// The Go stdlib time.Duration is an int64 number of nanoseconds.
func handle_stdlibTimeDuration(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {
	return handle_factor_with(area, rv, tag, label, []factor{
		{int64(time.Nanosecond), "nsec"},  // x1
		{int64(time.Microsecond), "μsec"}, // x1000
		{int64(time.Millisecond), "ms"},   // x1000 x1000
		{int64(time.Second), "seconds"},   // x1000 x1000 x1000
		{int64(time.Minute), "minutes"},   // x1000 x1000 x1000 x60
		{int64(time.Hour), "hours"},       // x1000 x1000 x1000 x60 x60
	})
}
