package autoconfig

type ConfigurableStruct interface{}

type AddressPort struct {
	Address string
	Port    int
}

type ExistingFile string

type ExistingDirectory string

type Password string

type EnumList int
