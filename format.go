package autoconfig

import (
	"fmt"
	"reflect"
	"strings"
)

// formatValue tries to format a plaintext summary of a reflect.Value.
func formatValue(rv *reflect.Value) string {
	kind := rv.Kind()

	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		return "Not configured"

	} else if stringer, ok := rv.Interface().(fmt.Stringer); ok { // n.b. matches if we have a T and (T) String() exists with value reciever
		return stringer.String()

	} else if rv.Kind() == reflect.String {
		return kind.String()

	} else if rv.CanInt() {
		return fmt.Sprintf("%d", rv.Int())

	} else if rv.CanUint() {
		return fmt.Sprintf("%d", rv.Uint())

	} else if rv.Kind() == reflect.Bool {
		return fmt.Sprintf("%v", rv.Bool())

	} else {

		if rv.CanAddr() {
			if stringer, ok := rv.Addr().Interface().(fmt.Stringer); ok { // n.b. matches if we have a T and (*T) String() exists with pointer reciever
				return stringer.String()
			}
		}

		// For a OneOf, we can try to use the OneOf's current value
		if rv.Kind() == reflect.Struct && rv.NumField() > 0 && rv.Field(0).Type() == reflect.TypeOf(OneOf("")) {
			if currentOneOf := rv.Field(0).String(); currentOneOf != "" {
				return currentOneOf
			}
		}

		if rv.Kind() == reflect.Pointer && !rv.IsNil() {
			// Pointer to something stringable?
			childItem := rv.Elem()
			childDisplayname := formatValue(&childItem)
			return "(" + childDisplayname + ")"
		}

		return "Configured"
	}
}

// formatLabel tries to generate a nice label from the automatic struct field.
func formatLabel(s string) string {
	// Mode: convert _ as spaces
	if strings.Contains(s, `_`) {
		return strings.TrimSpace(strings.ReplaceAll(s, `_`, ` `))
	}

	// Mode: convert CamelCase to spaces
	// Handle embedded acronyms (e.g. "TLSConfig" -> "TLS Config")
	var ret string
	var hold string // single previous uppercase character
	for i, ch := range s {
		_ = i
		if string(ch) == strings.ToUpper(string(ch)) {
			// uppercase
			if len(hold) == 0 && len(ret) > 0 && ret[len(ret)-1] != ' ' {
				// first uppercase after a previous lowercase = new word
				ret += " "
			}

			ret += hold
			hold = string(ch)
		} else {
			// lowercase
			if len(hold) > 0 {
				// first lowercase after a previous uppercase
				if len(ret) > 0 && ret[len(ret)-1] != ' ' {
					ret += " "
				}
				ret += hold
				hold = ""
			}
			ret += string(ch)
		}
	}
	if len(hold) > 0 {
		ret += hold
	}
	return ret
}
