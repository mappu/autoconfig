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
		A_File         ExistingFile
		A_Dir          ExistingDirectory
		Hostname       AddressPort
		Readonly       bool
		Something_Else *testInnerStruct
	}

	var myVar testStruct

	fmt.Printf("before = %#v\n", myVar)

	OpenDialog(&myVar, nil, "test dialog", func() {
		fmt.Printf("after  = %#v\n", myVar)
	})

	qt.QApplication_Exec()

}
