package autoconfig

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
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

type testIntegerBounds struct {
	Int8Min  int8
	Int8Max  int8
	Int16Min int16 `yprefix:"zim " ysuffix:" zam"`
	Int16Max int16
	Int32Min int32
	Int32Max int32
	Int64Min int64 `yprefix:"zim " ysuffix:" zam"`
	Int64Max int64

	UInt8Max  uint8
	UInt16Max uint16
	UInt32Max uint32
	UInt64Max uint64
}

func (t *testIntegerBounds) Reset() {
	t.Int8Min = math.MinInt8
	t.Int8Max = math.MaxInt8
	t.Int16Min = math.MinInt16
	t.Int16Max = math.MaxInt16
	t.Int32Min = math.MinInt32
	t.Int32Max = math.MaxInt32
	t.Int64Min = math.MinInt64
	t.Int64Max = math.MaxInt64

	t.UInt8Max = math.MaxUint8
	t.UInt16Max = math.MaxUint16
	t.UInt32Max = math.MaxUint32
	t.UInt64Max = math.MaxUint64
}

type testPrimitives struct {
	String  string
	Boolean bool
	Byte    byte
	Rune    rune
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

	// Put in a separate struct so they can be nilled on return (not JSON marshallable)
	ComplexPrimitives *struct {
		Complex64  complex64
		Complex128 complex128
	}

	ByteSlice []byte
}

func (t *testPrimitives) Reset() {
	t.String = fmt.Sprintf("testPrimitives.Reset() called at %s", time.Now().Format(time.RFC3339))
}

type testOneOf struct {
	SelectedType OneOf
	File         *ExistingFile      `yicon:"document-open"`
	Dir          *ExistingDirectory `yicon:"folder-open" ylabel:"Directory Custom Label"`
	Stdlib       *testStdlibTypes
}

type testTabGroup struct {
	_               TabGroup
	File            ExistingFile `yicon:"document-open"`
	FilePointer     *ExistingFile
	MultiLineString MultiLineString
	Struct          testOneOf
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

type testHijackedTypes struct {
	OrdinaryString         string
	OrdinaryStringDir      string
	OrdinaryStringPass     string
	OrdinaryStringPassword string

	OrdinaryInt         int
	OrdinaryIntDir      int
	OrdinaryIntPass     int
	OrdinaryIntPassword int
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
	FixedSizeArray    [4]string
	Deep_Pointer      *****TestInnerStruct
	H1                Header `ylabel:"Struct by value:"`
	DirectChild       TestInnerStruct
	H2                Header `ylabel:"Directly embedded struct:"`
	TestInnerStruct
}

type testMapTypes struct {
	Map_String_Empty   map[string]struct{}
	Map_String_String  map[string]string
	Map_String_Struct  map[string]TestInnerStruct
	Map_String_Pointer map[string]*TestInnerStruct
	Map_Int_String     map[int64]string

	PointerKeys *struct {
		Map_Pointer_String map[*TestInnerStruct]string
	}
}

type testStruct struct {
	H1              Header `ylabel:"This is the autoconfig test app"`
	Primitive_Types *testPrimitives
	Integer_Bounds  *testIntegerBounds
	Stdlib_Types    *testStdlibTypes
	Custom_Types    *testCustomTypes
	Hijack_Types    *testHijackedTypes
	Container_Types *testContainerTypes
	Map_Types       *testMapTypes
	OneOf           *testOneOf
	TabGroup        *testTabGroup
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
		OneOf: &testOneOf{
			SelectedType: "Stdlib", // Forcing the type by default should still allow changing it
		},
	}

	jbb, err := json.MarshalIndent(myVar, "", " ")
	if err != nil {
		t.Fatalf("Failed to JSON marshal old struct: %v", err)
	}

	fmt.Printf("BEFORE\n======\n\n%s\n\n", string(jbb))

	OpenDialog(&myVar, nil, "test dialog", func() {

		// Nil out some things that are interesting to test, but, not JSON
		// marshallable
		if myVar.Primitive_Types != nil {
			myVar.Primitive_Types.ComplexPrimitives = nil
		}
		if myVar.Map_Types != nil {
			myVar.Map_Types.PointerKeys = nil
		}

		jbb, err := json.MarshalIndent(myVar, "", " ")
		if err != nil {
			t.Fatalf("Failed to JSON marshal new struct: %v", err)
		}
		fmt.Printf("AFTER\n=====\n\n%s\n", string(jbb))
	})

	qt.QApplication_Exec()

}
