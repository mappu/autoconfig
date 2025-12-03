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
	- int
		- including uintptr, uint, and fixed-width versions
	- float
	- pointer (optional)
		- struct tags on the pointer are passed in to the child renderer
	- slice
	- struct
		- child structs by value, and embedded structs, are rendered inline
		- struct tags on the slice are passed in to each child renderer
	- empty struct
- Standard library types
	- time.Time
- Custom types
	- AddressPort
	- EnumList
	- ExistingDirectory
	- ExistingFile
	- Header
	- MultilineString
	- OneOf
	- Password
	- Any custom type that implements the `Renderer` interface

## Customization

Add struct tags to individual fields to customize the rendering:

|Tag      |Behaviour
|---------|------
|`ylabel` |Override label. If not present, the default label is the struct field's name with underscores replaced by spaces.
|`yenum`  |For "EnumList"; list of dropdown options, separated by double-semicolon (`;;`)
|`yfilter`|For "ExistingFile"; filter to apply in popup dialog
|`yicon`  |For "OneOf"; icon (either from theme, or with `:/` prefix for resource icon)

Implement these interfaces to customize the rendering:

|Interface       |Behaviour
|----------------|---------
|`Resetter`      |May be used with pointer receiver to reset your type to default values, if autoconfig constructed a new version of your type (used by OneOf, pointer, and slice)
|`Renderer`      |Add a fully custom Qt widget. Use with either value or pointer receiver.
|`fmt.Stringer`  |May be used to format some types for display

## Changelog

2025-12-03 v0.4.1

- Reduce flicker on Windows by enforcing Qt parent relationship during creation

2025-12-03 v0.4.0

- Reduce flicker on Windows by setting Qt UpdatesEnabled false during struct calculation
- OneOf: support Reset()
- Labels: Support String() on ExistingFile and ExistingDirectory
- Labels: Support labels on pointer types, with parenthesis
- Labels: Support labels on structs using OneOf
- Labels: Automatically convert CamelCase struct names to separated words
- Struct tags are now passed through to slice and pointer children (e.g. allowing `yfilter` on `*ExistingFile`)
- Use accurate type comparison for `time.Time` and `OneOf`, in case of naming conflict in other packages

2025-11-26 v0.3.0

- BREAKING: Rename `InitDefaulter` to `Resetter`, rename `Autoconfiger` to `Renderer`
- Add `OneOf`
- Renderer interface now supports being implemented on either value or pointer receiver
- `AddressPort` now renders a string description when used in a slice

2025-11-17 v0.2.0

- Support arbitrary pointers, slices, int, uintptr, float, `time.Time`
- Add `Header`
- Skip over unsupported types (func, interface, `unsafe.Pointer`)
- Fix cosmetic inconsistency when editing types

2025-11-15 v0.1.0

- Initial public release
