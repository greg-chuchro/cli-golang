package class_member_stores

import (
	"github.com/greg-chuchro/cli-golang/data_access"
)

// ClassMemberStoreFactory creates stores that access class members (fields and properties).
type ClassMemberStoreFactory struct {
	*ClassMemberStoreFactoryBase
}

func NewClassMemberStoreFactory() *ClassMemberStoreFactory {
	f := &ClassMemberStoreFactory{}
	f.ClassMemberStoreFactoryBase = NewClassMemberStoreFactoryBase(f)
	return f
}

func (this *ClassMemberStoreFactory) CreateFieldStore(obj any) data_access.IKeyValueStore[string, any] {
	return NewFieldStore(obj, this.BindingFlags)
}

func (this *ClassMemberStoreFactory) CreatePropertyStore(obj any) data_access.IKeyValueStore[string, any] {
	return NewPropertyStore(obj, this.BindingFlags)
}

var _ IClassMemberStoreFactory = (*ClassMemberStoreFactory)(nil)
