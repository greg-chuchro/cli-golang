package interface_components

import "reflect"

// Option represents a CLI option (field or property) with its metadata.
type Option struct {
	IsMultiValue bool
	Keys         []string
	MemberInfo   reflect.StructField
	ValueType    reflect.Type
	Value        string
}

// Parameter represents a CLI parameter for a subcommand with its metadata.
type Parameter struct {
	HasDefaultValue bool
	IsMultiValue    bool
	Keys            []string
	ParameterInfo   reflect.StructField
	ValueType       reflect.Type
	Value           string
}

// Subcommand represents a CLI subcommand (method) with its metadata.
type Subcommand struct {
	Keys       []string
	MethodInfo reflect.Method
	Parameters []Parameter
	ReturnType reflect.Type
}

// Submodule represents a CLI submodule (nested object) with its metadata.
type Submodule struct {
	Keys       []string
	MemberInfo reflect.StructField
}
