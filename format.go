package autoconfig

import (
	"fmt"
	"reflect"
)

// formatValue tries to format a plaintext summary of a reflect.Value.
func formatValue(rv *reflect.Value) string {
	kind := rv.Kind()

	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		return "Not configured"

	} else if stringer, ok := rv.Interface().(fmt.Stringer); ok { // n.b. matches if we have a T and (T) String() exists with value reciever
		return stringer.String()

	} else if stringer, ok := rv.Addr().Interface().(fmt.Stringer); ok { // n.b. matches if we have a T and (*T) String() exists with pointer reciever
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
		return "Object (" + rv.String() + ")"
	}
}
