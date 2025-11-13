package autoconfig

import (
	"fmt"
	"testing"

	qt "github.com/mappu/miqt/qt6"
)

func TestAutoConfig(t *testing.T) {

	qt.NewQApplication([]string{"test"})

	type testInnerStruct struct {
		Bar bool
	}

	type testStruct struct {
		Foo            string
		Something_Else *testInnerStruct
	}

	var myVar testStruct

	OpenDialog(&myVar, nil, "test dialog", func(result ConfigurableStruct) {
		fmt.Printf("%#v\n", result)
	})

	qt.QApplication_Exec()

}
