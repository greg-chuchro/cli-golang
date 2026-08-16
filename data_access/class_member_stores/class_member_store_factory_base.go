package class_member_stores

import (
	"github.com/greg-chuchro/cli-golang/data_access"
	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// IClassMemberStoreFactory creates key-value stores that access class members (fields and properties).
type IClassMemberStoreFactory interface {
	CreateClassMemberStore(obj any) data_access.IKeyValueStore[string, any]
	CreateFieldStore(obj any) data_access.IKeyValueStore[string, any]
	CreatePropertyStore(obj any) data_access.IKeyValueStore[string, any]
}

// ClassMemberStoreFactoryBase is the base factory for creating stores that access class members.
// It follows the guideline delegation pattern: the concrete factory passes itself (factory) into the base
// so the base invokes the concrete's overridden CreateFieldStore/CreatePropertyStore through the interface.
type ClassMemberStoreFactoryBase struct {
	factory          IClassMemberStoreFactory
	AccessFields     bool
	AccessProperties bool
	BindingFlags     BindingFlags
}

func NewClassMemberStoreFactoryBase(factory IClassMemberStoreFactory) *ClassMemberStoreFactoryBase {
	return &ClassMemberStoreFactoryBase{
		factory:          factory,
		AccessFields:     false,
		AccessProperties: true,
		BindingFlags:     DefaultLookup,
	}
}

func (this *ClassMemberStoreFactoryBase) CreateClassMemberStore(obj any) data_access.IKeyValueStore[string, any] {
	if this.AccessFields && this.AccessProperties {
		return this.CreateFieldAndPropertyStore(obj)
	}
	if this.AccessFields {
		return this.factory.CreateFieldStore(obj)
	}
	if this.AccessProperties {
		return this.factory.CreatePropertyStore(obj)
	}
	panic(exceptions.NoAccessConfigured())
}

func (this *ClassMemberStoreFactoryBase) CreateFieldAndPropertyStore(obj any) data_access.IKeyValueStore[string, any] {
	return data_access.NewDualKeyValueStore[string, any](this.factory.CreateFieldStore(obj), this.factory.CreatePropertyStore(obj))
}

func (this *ClassMemberStoreFactoryBase) CreateFieldStore(obj any) data_access.IKeyValueStore[string, any] {
	return this.factory.CreateFieldStore(obj)
}

func (this *ClassMemberStoreFactoryBase) CreatePropertyStore(obj any) data_access.IKeyValueStore[string, any] {
	return this.factory.CreatePropertyStore(obj)
}
