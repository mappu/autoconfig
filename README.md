# autoconfig

Autoconfig allows you to edit any Go struct with a Qt interface.

```
struct Foo {                     []----------------------[]
    Name string                  |  Name:  [___________]  |
}                                |                [Save]  |
                                 []----------------------[]
```

## Supported types

- string, bool, AddressPort, ExistingFile, Password

It supports the struct tags:

- `ylabel` - Override label. Otherwise, use field name
- `yfilter` - For "ExistingFile"; filter to apply in popup dialog
- `yport` - For ”AddressPort"
