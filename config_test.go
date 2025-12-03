package autoconfig

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	qt "github.com/mappu/miqt/qt6"
)

type TestInnerStruct struct {
	Bar bool
}

func (t *TestInnerStruct) String() string {
	return fmt.Sprintf("my bar is %v", t.Bar)
}

type testPrimitives struct {
	String  string
	Boolean bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	UInt    uint
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	Uintptr uintptr
	Float32 float32
	Float64 float64
}

type testOneOf struct {
	SelectedType OneOf
	File         *ExistingFile      `yicon:"document-open"`
	Dir          *ExistingDirectory `yicon:"folder-open" ylabel:"Directory Custom Label"`
	Stdlib       *testStdlibTypes
}

type testCustomTypes struct {
	H1             Header       `ylabel:"Types by value"`
	A_File         ExistingFile `yfilter:"Text files (*.txt);;All files (*)"`
	A_Dir          ExistingDirectory
	Hostname       AddressPort
	Multiple_Lines MultiLineString
	FooPassword    Password

	H2                 Header        `ylabel:"Types by pointer"`
	A_File_Ptr         *ExistingFile `yfilter:"Text files (*.txt);;All files (*)"`
	A_Dir_Ptr          *ExistingDirectory
	Hostname_Ptr       *AddressPort
	Multiple_Lines_Ptr *MultiLineString
	FooPassword_Ptr    *Password
}

type testStdlibTypes struct {
	Time      time.Time
	TLSConfig *tls.Config
	NetDialer *net.Dialer
}

type testContainerTypes struct {
	EmptyStruct       struct{}
	Empty_By_Pointer  *struct{}
	Struct_By_Pointer *TestInnerStruct
	Struct_By_Slice   []TestInnerStruct
	Custom_By_Slice   []ExistingFile `ylabel:"Custom by slice (ylabel)" yfilter:"Text files (*.txt);;All files (*)"` // Tag attributes on slices are passed into the child
	Struct_Ptr_Slice  []*TestInnerStruct
	Deep_Pointer      *****TestInnerStruct
	H1                Header `ylabel:"Struct by value:"`
	DirectChild       TestInnerStruct
	H2                Header `ylabel:"Directly embedded struct:"`
	TestInnerStruct
}

type testStruct struct {
	H1              Header `ylabel:"This is the autoconfig test app"`
	Primitive_Types *testPrimitives
	Stdlib_Types    *testStdlibTypes
	Custom_Types    *testCustomTypes
	Container_Types *testContainerTypes
	OneOf           *testOneOf
}

func TestAutoConfig(t *testing.T) {

	qt.NewQApplication([]string{"test"})

	myVar := testStruct{
		Stdlib_Types: &testStdlibTypes{
			Time: time.Now(),
		},
		Container_Types: &testContainerTypes{
			Struct_By_Slice: []TestInnerStruct{
				TestInnerStruct{Bar: true},
				TestInnerStruct{Bar: false},
			},
			Struct_Ptr_Slice: []*TestInnerStruct{
				&TestInnerStruct{Bar: true},
			},
		},
	}

	jbb, _ := json.MarshalIndent(myVar, "", " ")
	fmt.Printf("BEFORE\n======\n\n%s\n\n", string(jbb))

	OpenDialog(&myVar, nil, "test dialog", func() {

		jbb, _ := json.MarshalIndent(myVar, "", " ")
		fmt.Printf("AFTER\n=====\n\n%s\n", string(jbb))
	})

	qt.QApplication_Exec()

}
