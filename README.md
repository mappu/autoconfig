# autoconfig

![](https://img.shields.io/badge/License-MIT-green)
[![Go Reference](https://pkg.go.dev/badge/github.com/mappu/autoconfig.svg)](https://pkg.go.dev/github.com/mappu/autoconfig)

Autoconfig allows you to edit any Go struct with a Qt interface [based on MIQT](https://github.com/mappu/miqt).

```
type Foo struct {                []----------------------[]
    Name string                  |  Name:  [___________]  |
}                                |                [Save]  |
                                 []----------------------[]
```

## Usage

Creating a dialog:

```golang
// Passed in struct should be a pointer value
var foo MyStruct
autoconfig.OpenDialog(&foo, nil, "Dialog title", func() {
	// The value of 'foo' has been updated
})
```

Embedding into an existing layout:

```golang
var foo MyStruct
saveCallback := autoconfig.MakeConfigArea(&foo, qt6.QFormLayout)

// To save changes from the GUI into the struct, call the saveCallback() function.
// However, warning that nested fields may be mutated automatically without calling.
```

Only public fields are supported. This is a limitation of the standard library `reflect` package.

## Supported types

- Primitive types
	- string
	- bool
	- int (including uintptr, uint, and fixed-width versions)
	- float
	- pointer (optional)
	- slice
	- struct
- Standard library types
	- time.Time
- Custom types
	- AddressPort
	- EnumList
	- ExistingDirectory
	- ExistingFile
	- Header
	- MultilineString
	- Password
	- Any custom type that implements the `Autoconfiger` interface

## Customization

Add struct tags to individual fields to customize the rendering:

|Tag      |Behaviour
|---------|------
|`ylabel` |Override label. If not present, the default label is the struct field's name with underscores replaced by spaces.
|`yfilter`|For "ExistingFile"; filter to apply in popup dialog
|`yenum`  |For "EnumList"; list of dropdown options, separated by double-semicolon (`;;`)

Implement these interfaces to customize the rendering:

|Interface       |Behaviour
|----------------|---------
|`Resetter`      |May be used with pointer receiver to reset your type to default values, if autoconfig constructed a new version of your type
|`Autoconfiger`  |Add a fully custom Qt widget
|`fmt.Stringer`  |May be used to format some types for display
