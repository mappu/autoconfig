package qspinbox

import (
	"strconv"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

// QUint64SpinBox is a QSpinBox that works on uint64 types.
type QUint64SpinBox struct {
	*qt.QAbstractSpinBox

	minimum, maximum, value uint64
	prefix, suffix          string
}

func (s *QUint64SpinBox) Minimum() uint64 {
	return s.minimum
}

func (s *QUint64SpinBox) SetMinimum(newMinimum uint64) {
	s.minimum = newMinimum
}

func (s *QUint64SpinBox) Maximum() uint64 {
	return s.maximum
}

func (s *QUint64SpinBox) SetMaximum(newMaximum uint64) {
	s.maximum = newMaximum
}

func (s *QUint64SpinBox) Value() uint64 {
	return s.value
}

func (s *QUint64SpinBox) SetValue(newValue uint64) {
	s.LineEdit().SetText(s.textFromValue(newValue))
	s.value = newValue
}

func (s *QUint64SpinBox) Prefix() string {
	return s.prefix
}

func (s *QUint64SpinBox) SetPrefix(newPrefix string) {
	s.prefix = newPrefix
	s.SetValue(s.value) // re-render
}

func (s *QUint64SpinBox) Suffix() string {
	return s.suffix
}

func (s *QUint64SpinBox) SetSuffix(newSuffix string) {
	s.suffix = newSuffix
	s.SetValue(s.value) // re-render
}

func (s *QUint64SpinBox) textFromValue(value uint64) string {
	return s.prefix + strconv.FormatUint(value, 10) + s.suffix
}

func (s *QUint64SpinBox) valueFromText(displayText string) (uint64, error) {

	// Strip suffix + prefix if we have them
	// Also, allow the user to paste in a number omitting the suffix/prefix

	if s.prefix != "" && strings.HasPrefix(displayText, s.prefix) {
		displayText = displayText[len(s.prefix):]
	}

	if s.suffix != "" && strings.HasSuffix(displayText, s.suffix) {
		displayText = displayText[0 : len(displayText)-len(s.suffix)]
	}

	// After stripping prefix+suffix, only the numeric value remains
	return strconv.ParseUint(displayText, 10, 64)
}

// NewQUint64SpinBox constructs a new QUint64SpinBox.
func NewQUint64SpinBox(parent *qt.QWidget) *QUint64SpinBox {
	s := &QUint64SpinBox{}

	if parent == nil {
		s.QAbstractSpinBox = qt.NewQAbstractSpinBox2()
	} else {
		s.QAbstractSpinBox = qt.NewQAbstractSpinBox(parent)
	}

	s.QAbstractSpinBox.OnStepBy(func(super func(steps int), steps int) {
		s.SetValue(AddSaturatingUnsigned(s.value, steps))
	})

	// By default, our widget size is 0 pixels wide(??)
	// FIXME this should properly be based on OS metrics and needs to account for the prefix/suffix length
	// QSpinBox::sizeHint() does this via fontMetrics()->horizontalAdvance(sample text)
	s.SetMinimumWidth(190)

	s.QAbstractSpinBox.OnValidate(func(super func(input string, pos *int) qt.QValidator__State, input string, pos *int) qt.QValidator__State {
		val, err := s.valueFromText(input)
		if err != nil {
			return qt.QValidator__Invalid
		}

		if val < s.minimum || val > s.maximum {
			return qt.QValidator__Invalid
		}

		return qt.QValidator__Acceptable
	})

	s.QAbstractSpinBox.OnStepEnabled(func(super func() qt.QAbstractSpinBox__StepEnabledFlag) qt.QAbstractSpinBox__StepEnabledFlag {
		if s.value == s.minimum {
			return qt.QAbstractSpinBox__StepUpEnabled
		} else if s.value == s.maximum {
			return qt.QAbstractSpinBox__StepDownEnabled
		} else {
			return qt.QAbstractSpinBox__StepUpEnabled | qt.QAbstractSpinBox__StepDownEnabled
		}
	})

	s.LineEdit().OnTextEdited(func(input string) {
		val, err := s.valueFromText(input)
		if err != nil {
			// Edited to an invalid value
			// Undo text change, revert to current value
			s.LineEdit().SetText(s.textFromValue(s.value))
			return
		}

		// Text changed, update our internal model
		s.value = val

		// Probably no need to change the text display, that's already done
		// Unless the user pasted in something without the prefix/suffix
		expect := s.textFromValue(s.value)
		if input != expect {
			s.LineEdit().SetText(expect) // n.b. may re-trigger OnTextEdited once
			s.LineEdit().SetCursorPosition(len(expect) - len(s.suffix))
		}
	})

	return s
}
