package interface_component_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/data_access/class_member_stores"
	cms "github.com/greg-chuchro/cli-golang/data_access/class_member_stores"
	dae "github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// CliComponentStoreFactory creates CLI component stores (options, subcommands, submodules) from class objects.
type CliComponentStoreFactory struct {
	AccessFields              bool
	AccessProperties          bool
	BindingFlags              cms.BindingFlags
	ClassMemberStringifier    data_access.IClassMemberStringifier
	FormatProvider            string
	MethodBindingFlags        cms.BindingFlags
	OptionAccessValidator     data_access.IAccessValidator
	SubcommandAccessValidator data_access.IAccessValidator
	SubmoduleAccessValidator  data_access.IAccessValidator
	CollectionStoreFactory    data_access.ICollectionElementStoreFactory
	CompositeValueConverter   data_access.ICompositeValueConverter[any]
	EnableShadowing           bool
}

// DefaultLookup matches C#'s DefaultLookup (Instance, Static, Public, IgnoreCase).
var DefaultLookup = cms.BindingFlags{Instance: true, Static: true, Public: true, IgnoreCase: true}

func NewCliComponentStoreFactory() *CliComponentStoreFactory {
	f := &CliComponentStoreFactory{
		AccessFields:           true,
		AccessProperties:       true,
		BindingFlags:           DefaultLookup,
		ClassMemberStringifier: &data_access.ClassMemberStringifier{},
		FormatProvider:         "R",
		MethodBindingFlags:     DefaultLookup,
		EnableShadowing:        false,
	}
	return f
}

func (this *CliComponentStoreFactory) compositeValueConverter() data_access.ICompositeValueConverter[any] {
	if this.CompositeValueConverter != nil {
		return this.CompositeValueConverter
	}
	base := &data_access.ValueConverter{FormatProvider: this.FormatProvider}
	this.CompositeValueConverter = data_access.NewCompositeValueConverter[any](base, this.collectionStoreFactory())
	return this.CompositeValueConverter
}

func (this *CliComponentStoreFactory) collectionStoreFactory() data_access.ICollectionElementStoreFactory {
	if this.CollectionStoreFactory != nil {
		return this.CollectionStoreFactory
	}
	this.CollectionStoreFactory = &data_access.CollectionElementStoreFactory{}
	return this.CollectionStoreFactory
}

func (this *CliComponentStoreFactory) optionAccessValidator() data_access.IAccessValidator {
	if this.OptionAccessValidator != nil {
		return this.OptionAccessValidator
	}
	this.OptionAccessValidator = data_access.NewOptionAccessValidator(this.compositeValueConverter())
	return this.OptionAccessValidator
}

func (this *CliComponentStoreFactory) subcommandAccessValidator() data_access.IAccessValidator {
	if this.SubcommandAccessValidator != nil {
		return this.SubcommandAccessValidator
	}
	this.SubcommandAccessValidator = &data_access.SubcommandAccessValidator{}
	return this.SubcommandAccessValidator
}

func (this *CliComponentStoreFactory) submoduleAccessValidator() data_access.IAccessValidator {
	if this.SubmoduleAccessValidator != nil {
		return this.SubmoduleAccessValidator
	}
	this.SubmoduleAccessValidator = data_access.NewSubmoduleAccessValidator(this.compositeValueConverter())
	return this.SubmoduleAccessValidator
}

func (this *CliComponentStoreFactory) CreateOptionStore(obj any) IOptionStore {
	var store data_access.ICliKeyValueStore[any]
	if this.AccessFields && this.AccessProperties {
		store = this.createFieldAndPropertyStore(obj, this.optionAccessValidator(), this.compositeValueConverter())
	} else if this.AccessFields {
		store = this.createFieldStore(obj, this.optionAccessValidator(), this.compositeValueConverter())
	} else if this.AccessProperties {
		store = this.createPropertyStore(obj, this.optionAccessValidator(), this.compositeValueConverter())
	} else {
		panic(dae.DataAccessErrors.NoAccessConfigured())
	}
	return NewOptionStore(store)
}

func (this *CliComponentStoreFactory) CreateSubcommandExecutor(method reflect.Method, obj any) ISubcommandExecutor {
	bound := reflect.Value{}
	if obj != nil {
		bound = reflect.ValueOf(obj).MethodByName(method.Name)
	} else {
		recvType := method.Type.In(0)
		var recv reflect.Value
		if recvType.Kind() == reflect.Ptr {
			recv = reflect.New(recvType.Elem())
		} else {
			recv = reflect.New(recvType)
		}
		bound = recv.MethodByName(method.Name)
	}
	executor := class_member_stores.NewCliMethodExecutor(bound, obj, this.MethodBindingFlags, this.ClassMemberStringifier, this.compositeValueConverter())
	return NewSubcommandExecutor(executor)
}

func (this *CliComponentStoreFactory) CreateSubcommandExecutorWithOptions(method reflect.Method, obj any) ISubcommandExecutorWithOptions {
	return NewSubcommandExecutorWithOptions(this.CreateSubcommandExecutor(method, obj), this.CreateOptionStore(obj))
}

func (this *CliComponentStoreFactory) CreateSubcommandStore(obj any) ISubcommandStore {
	return NewSubcommandStore(class_member_stores.NewMethodInfoStore(obj, this.MethodBindingFlags, this.ClassMemberStringifier, this.subcommandAccessValidator()))
}

func (this *CliComponentStoreFactory) CreateSubmoduleStore(obj any) ISubmoduleStore {
	converter := &data_access.ReadOnlyPassThroughConverter[any]{}
	var store data_access.ICliKeyValueStore[any]
	if this.AccessFields && this.AccessProperties {
		store = this.createFieldAndPropertyStore(obj, this.submoduleAccessValidator(), converter)
	} else if this.AccessFields {
		store = this.createFieldStore(obj, this.submoduleAccessValidator(), converter)
	} else if this.AccessProperties {
		store = this.createPropertyStore(obj, this.submoduleAccessValidator(), converter)
	} else {
		panic(dae.DataAccessErrors.NoAccessConfigured())
	}
	return NewSubmoduleStore(store)
}

func (this *CliComponentStoreFactory) createFieldAndPropertyStore(obj any, validator data_access.IAccessValidator, converter data_access.ICompositeValueConverter[any]) data_access.ICliKeyValueStore[any] {
	return &cliDualKeyValueStore{
		primary:   this.createFieldStore(obj, validator, converter),
		secondary: this.createPropertyStore(obj, validator, converter),
	}
}

func (this *CliComponentStoreFactory) createFieldStore(obj any, validator data_access.IAccessValidator, converter data_access.ICompositeValueConverter[any]) data_access.ICliKeyValueStore[any] {
	return class_member_stores.NewCliFieldStore(obj, this.BindingFlags, this.ClassMemberStringifier, validator, converter)
}

func (this *CliComponentStoreFactory) createPropertyStore(obj any, validator data_access.IAccessValidator, converter data_access.ICompositeValueConverter[any]) data_access.ICliKeyValueStore[any] {
	return class_member_stores.NewCliPropertyStore(obj, this.BindingFlags, this.ClassMemberStringifier, validator, converter)
}
