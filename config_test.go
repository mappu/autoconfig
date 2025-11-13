package autoconfig

import (
	"fmt"
	"testing"

	qt "github.com/mappu/miqt/qt6"
)

func TestAutoConfig(t *testing.T) {

	qt.NewQApplication([]string{"test"})

	type testStruct struct {
		Foo string
	}

	var myVar testStruct

	OpenDialog(&myVar, nil, "test dialog", func(result ConfigurableStruct) {
		fmt.Printf("%#v\n", result)
	})

	qt.QApplication_Exec()

}
