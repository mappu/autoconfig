package autoconfig

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	qt "github.com/mappu/miqt/qt6"
)

type testInnerStruct struct {
	Bar bool
}

func (t *testInnerStruct) String() string {
	return fmt.Sprintf("my bar is %v", t.Bar)
}

type testStruct struct {
	Foo               string
	A_File            ExistingFile
	A_Dir             ExistingDirectory
	Hostname          AddressPort
	Multiple_Lines    MultiLineString
	H1                Header `ylabel:"Optional settings with long label spanning columns"`
	FooPassword       Password
	Readonly          bool
	Struct_By_Pointer *testInnerStruct
	Struct_By_Slice   []testInnerStruct
	Time              time.Time
	TLSConfig         *tls.Config
}

func TestAutoConfig(t *testing.T) {

	qt.NewQApplication([]string{"test"})

	var myVar testStruct
	myVar.Time = time.Now()
	myVar.Struct_By_Slice = []testInnerStruct{
		testInnerStruct{Bar: true},
		testInnerStruct{Bar: false},
	}

	fmt.Printf("before = %#v\n", myVar)

	OpenDialog(&myVar, nil, "test dialog", func() {
		fmt.Printf("after  = %#v\n", myVar)
	})

	qt.QApplication_Exec()

}
