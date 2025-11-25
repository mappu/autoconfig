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

## Supported types

- string, bool, `*struct`
- Any custom types (many types included in package)

## Customization

Struct tags:

- `ylabel` - Override label. If not present, the default label is the struct field's name with underscores replaced by spaces.
- `yfilter` - For "ExistingFile"; filter to apply in popup dialog
- `yenum` - For "EnumList"; list of dropdown options, separated by double-semicolon (`;;`)

Interfaces:

- `InitDefaulter` - May be used if autoconfig needs to construct a new version of your type
- `Autoconfiger` - Add a fully custom Qt widget
- `fmt.Stringer` - May be used to format some types for display

## Notes

- Passed in struct should be a pointer value
- Call the saver, but, warning that some fields may be mutated automatically without calling
- Public fields only
